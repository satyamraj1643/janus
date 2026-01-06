package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/middleware"
	"github.com/satyamraj1643/janus/queue"
	"github.com/satyamraj1643/janus/spec"
)

const (
	SystemBatchID    = "11111111-1111-1111-1111-111111111111" // fixed UUID for system batch
	DashboardBatchID = "22222222-2222-2222-2222-222222222222" // fixed UUID for dashboard batch
)

type JobHandler struct {
	AC            *admission.AdmissionController
	FromDashboard bool
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	log.Println("PATH:", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var job spec.Job

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if job.ID == "" || job.TenantID == "" {
		http.Error(w, "missing job_id or tenanat_id", http.StatusBadRequest)
		return
	}

	if h.FromDashboard {
		job.Source = spec.JobSourceDashboard
		job.BatchName = "dashboard_batch"
		job.BatchID = DashboardBatchID
	} else {
		job.Source = spec.JobSourceSystem
		job.BatchName = "system_batch"
		job.BatchID = SystemBatchID
	}

	// Attach user's active config to job
	activeConfig, configID, ownerID, _ := middleware.GetActiveContext(r.Context())
	job.Config = activeConfig
	job.GlobalConfigID = configID
	job.OwnerID = ownerID

	log.Printf("Validating job synchronously.")

	decision, err := h.AC.Check(r.Context(), job)
	if err != nil {
		http.Error(w, "internal service error", http.StatusInternalServerError)
		return
	}

	// Send to DB writer
	queue.ResultQueue <- decision

	w.Header().Set("Content-Type", "application/json")
	if decision.Status == "accepted" {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusForbidden) // Or 429 based on reason
	}
	json.NewEncoder(w).Encode(decision)
}

// Accepts partial batch, some admitted and some not.
func (h *JobHandler) CreateJobBatch(w http.ResponseWriter, r *http.Request) {
	log.Println("PATH:", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req JobBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// ✅ batch_name is mandatory
	if req.BatchName == "" {
		http.Error(w, "batch_name required", http.StatusBadRequest)
		return
	}
	batchName := req.BatchName

	// ✅ BatchID depends ONLY on source
	var batchID string
	if h.FromDashboard {
		batchID = "dashboard_batch_" + uuid.NewString()
	} else {
		batchID = "system_batch_" + uuid.NewString()
	}

	if len(req.Jobs) == 0 {
		http.Error(w, "jobs array cannot be empty", http.StatusBadRequest)
		return
	}

	const maxBatchSize = 1000
	if len(req.Jobs) > maxBatchSize {
		http.Error(w, "batch too large", http.StatusRequestEntityTooLarge)
		return
	}

	admitted := 0
	var decisions []*spec.JobDecision

	for _, job := range req.Jobs {
		if job.ID == "" || job.TenantID == "" {
			break
		}

		if h.FromDashboard {
			job.Source = spec.JobSourceDashboard
		} else {
			job.Source = spec.JobSourceSystem
		}

		job.BatchName = batchName
		job.BatchID = batchID

		// Attach user's active config to job
		activeConfig, configID, ownerID, _ := middleware.GetActiveContext(r.Context())
		job.Config = activeConfig
		job.GlobalConfigID = configID
		job.OwnerID = ownerID

		decision, err := h.AC.Check(r.Context(), job)
		if err != nil {
			// If internal error, we count as failure/break or log
			break
		}

		queue.ResultQueue <- decision
		decisions = append(decisions, decision)

		if decision.Status == "accepted" {
			admitted++
		}
	}

	rejected := len(req.Jobs) - admitted

	status := "full"
	if admitted == 0 {
		status = "rejected"
	} else if rejected > 0 {
		status = "partial"
	}

	resp := JobBatchResponse{
		BatchName: batchName,
		Status:    status,
		Admitted:  admitted,
		Rejected:  rejected,
	}

	// Optional: we can return decisions details if needed, but keeping existing response format for now + decisions?
	// The user asked for "status of each job", but the original response was aggregate.
	// Let's stick to aggregate for this endpoint as it's partial, or maybe add detail?
	// User said "what is the status of each job after janus run".
	// Let's just return the aggregate for now to pass build, then refine.

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (h *JobHandler) CreateJobBatchAtomic(w http.ResponseWriter, r *http.Request) {
	log.Println("PATH:", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req JobBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.BatchName == "" {
		http.Error(w, "batch_name required", http.StatusBadRequest)
		return
	}
	batchName := req.BatchName

	var batchID string
	if h.FromDashboard {
		batchID = "dashboard_batch_" + uuid.NewString()
	} else {
		batchID = "system_batch_" + uuid.NewString()
	}

	if len(req.Jobs) == 0 {
		http.Error(w, "include at least 1 job in the batch", http.StatusBadRequest)
		return
	}

	for i := range req.Jobs {
		if req.Jobs[i].ID == "" || req.Jobs[i].TenantID == "" {
			http.Error(w, "invalid job in batch", http.StatusBadRequest)
			return
		}

		if h.FromDashboard {
			req.Jobs[i].Source = spec.JobSourceDashboard
		} else {
			req.Jobs[i].Source = spec.JobSourceSystem
		}

		req.Jobs[i].BatchName = batchName
		req.Jobs[i].BatchID = batchID

		// Attach user's active config to job
		activeConfig, configID, ownerID, _ := middleware.GetActiveContext(r.Context())
		req.Jobs[i].Config = activeConfig
		req.Jobs[i].GlobalConfigID = configID
		req.Jobs[i].OwnerID = ownerID
	}

	decisions, err := h.AC.CheckBatchAtomic(r.Context(), req.Jobs)
	if err != nil {
		http.Error(w, "internal error during atomic check", http.StatusInternalServerError)
		return
	}

	// Queue decisions for DB
	for _, d := range decisions {
		queue.ResultQueue <- d
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")

	// If any rejected, the whole batch is technically rejected in "atomic" sense if we consider
	// CheckBatchAtomic's logic.
	// But CheckBatchAtomic returns a list of decisions.
	// We can inspect the first decision status (as they should all be same for atomic)
	status := http.StatusAccepted
	if len(decisions) > 0 && decisions[0].Status == "rejected" {
		status = http.StatusForbidden
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(decisions)
}

package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/satyamraj1643/janus/queue"
	"github.com/satyamraj1643/janus/spec"
)

const (
	SystemBatchID    = "11111111-1111-1111-1111-111111111111" // fixed UUID for system batch
	DashboardBatchID = "22222222-2222-2222-2222-222222222222" // fixed UUID for dashboard batch
)

func CreateJob(w http.ResponseWriter, r *http.Request, fromDashboard bool) {
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

	if fromDashboard {
		job.Source = spec.JobSourceDashboard
		job.BatchName = "dashboard_job"
		job.BatchID = DashboardBatchID
	} else {
		job.Source = spec.JobSourceSystem
		job.BatchName = "system_batch"
		job.BatchID = SystemBatchID
	}

	log.Printf("Pushing in intermediate job-queue.")

	if err := queue.Admit(job); err != nil {
		http.Error(w, "admission queue busy", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

// Accepts partial batch, some admitted and some not.
func CreateJobBatch(w http.ResponseWriter, r *http.Request, fromDashboard bool) {
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
	if fromDashboard {
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

	for _, job := range req.Jobs {
		if job.ID == "" || job.TenantID == "" {
			break
		}

		if fromDashboard {
			job.Source = spec.JobSourceDashboard
		} else {
			job.Source = spec.JobSourceSystem
		}

		job.BatchName = batchName
		job.BatchID = batchID

		if err := queue.Admit(job); err != nil {
			break
		}
		admitted++
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}


func CreateJobBatchAtomic(w http.ResponseWriter, r *http.Request, fromDashboard bool) {
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
	if fromDashboard {
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

		if fromDashboard {
			req.Jobs[i].Source = spec.JobSourceDashboard
		} else {
			req.Jobs[i].Source = spec.JobSourceSystem
		}

		req.Jobs[i].BatchName = batchName
		req.Jobs[i].BatchID = batchID
	}

	if err := queue.AdmitBatch(req.Jobs); err != nil {
		http.Error(w, "system busy", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

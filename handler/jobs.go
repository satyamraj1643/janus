package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/satyamraj1643/janus/queue"
	"github.com/satyamraj1643/janus/spec"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Pushing in intermediate job-queue.")

	if err := queue.Admit(job); err != nil {
		http.Error(w, "admission queue busy", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

// Accepts partial batch, some admitted and some not.
func CreateJobBatch(w http.ResponseWriter, r *http.Request) {
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
		BatchName: req.BatchName,
		Status:    status,
		Admitted:  admitted,
		Rejected:  rejected,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func CreateJobBatchAtomic(w http.ResponseWriter, r *http.Request) {
	log.Println("PATH:", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	defer r.Body.Close()

	var req JobBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
	}

	if req.BatchName == "" || len(req.Jobs) == 0 {
		http.Error(w, "invalid batch", 400)
		return
	}

	for _, job := range req.Jobs {
		if job.ID == "" || job.TenantID == "" {
			http.Error(w, "invalid job in batch", 400)
			return
		}
	}

	if err := queue.AdmitBatch(req.Jobs); err != nil {
		http.Error(w, "system busy", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

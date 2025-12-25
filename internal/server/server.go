package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/spec"
)

type Server struct {
	ac *admission.AdmissionController
}

func NewServer(ac *admission.AdmissionController) *Server {
	return &Server{ac: ac}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("GET /stats", s.handleStats)
	mux.HandleFunc("GET /config", s.handleConfig)
	mux.HandleFunc("POST /simulate", s.handleSimulate)

	// Enable CORS for Wails
	handler := middlewareCors(mux)

	fmt.Printf("Server listening on %s\n", addr)
	return http.ListenAndServe(addr, handler)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	// Requires GetStats() in controller.go
	stats := s.ac.GetStats()
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var newPolicy policy.Policy
		if err := json.NewDecoder(r.Body).Decode(&newPolicy); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.ac.UpdatePolicy(&newPolicy)
		w.Write([]byte("policy updated"))
		return
	}
	json.NewEncoder(w).Encode(s.ac.Policy)
}

func (s *Server) handleSimulate(w http.ResponseWriter, r *http.Request) {
	var job spec.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Auto-generate ID if missing
	if job.ID == "" {
		job.ID = fmt.Sprintf("sim_%d", time.Now().UnixNano())
	}

	// Run Admission Check
	start := time.Now()
	err := s.ac.Check(r.Context(), job)
	duration := time.Since(start)

	resp := map[string]interface{}{
		"admitted":   err == nil,
		"latency_ms": duration.Milliseconds(),
		"error":      "",
	}
	if err != nil {
		resp["error"] = err.Error()
	}

	json.NewEncoder(w).Encode(resp)
}

// CORS Middleware to allow Wails frontend to connect
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

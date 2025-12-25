package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/internal/server"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/spec"
)

func main() {
	// 1. Load Policy
	p, err := policy.LoadPolicy("config/janus.json")
	if err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}
	fmt.Printf("loaded policy v%d with %d dependencies\n", p.Version, len(p.Dependencies))

	// 2. Load Jobs (Optional, just ensuring valid jobs.json)
	_, err = loadJobs("config/jobs.json")
	if err != nil {
		log.Printf("Warning: Failed to load jobs.json: %v", err)
	}

	// 3. Connect to Redis (Assuming localhost:6379 for now)
	redisStore := store.NewRedisStore("localhost:6379")
	if err := redisStore.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✅ REDIS CONNECTED")

	// 4. Initialize Admission Controller
	ac := admission.NewAdmissionController(p, redisStore)

	// 5. Start HTTP Server
	srv := server.NewServer(ac)
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loadJobs(path string) ([]spec.Job, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var jobs []spec.Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

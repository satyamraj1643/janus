package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	// "time"

	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/spec"
)

func main() {
	// 1. Load Policy
	p, err := policy.LoadPolicy("config/janus.json")
	if err != nil {
		log.Fatalf("FAILED to load policy: %v", err)
	}
	fmt.Printf("loaded policy v%d with %d dependencies\n", p.Version, len(p.Dependencies))

	// 2. Initialize Store (Redis)
	// assuming localhost:6379 for now
	redisStore := store.NewRedisStore("localhost:6379")
	ctx := context.Background()
	if err := redisStore.Ping(ctx); err != nil {
		log.Fatalf("FAILED to connect to Redis: %v", err)
	}
	fmt.Println("✅ REDIS CONNECTED")

	// CLEANUP: Reset Redis state for this simulation run
	if err := redisStore.Flush(ctx); err != nil {
		log.Fatalf("FAILED to flush Redis: %v", err)
	}

	fmt.Println("🧹 REDIS FLUSHED")

	// 3. Initialize Admission Controller
	ac := admission.NewAdmissionController(p, redisStore)

	// 4. Load Sample Jobs
	jobs, err := loadJobs("config/jobs.json")
	if err != nil {
		log.Fatalf("FAILED to load jobs: %v", err)
	}
	fmt.Printf("loaded %d sample jobs\n", len(jobs))

	// 5. Run Simulation
	fmt.Println("\n--- STARTING ADMISSION SIMULATION ---")

	for _, job := range jobs {
		fmt.Printf("Processing Job %s (Tenant: %s)... ", job.ID, job.TenantID)

		err := ac.Check(ctx, job)

		// if jobNumber == 3 {
		// 	fmt.Println("Waiting for window to expire...")
        //     time.Sleep(2 * time.Second) 
		// }


		if err != nil {
			fmt.Printf("REJECTED: %v\n", err)
		} else {
			fmt.Printf("ADMITTED\n")
		}
	}
	fmt.Println("--- SIMULATION COMPLETE ---")
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

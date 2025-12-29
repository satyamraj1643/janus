package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/satyamraj1643/janus/db"
	"github.com/satyamraj1643/janus/handler"
	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/middleware"
	"github.com/satyamraj1643/janus/worker"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Initialise admission controller 
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR not set")
	}

	redisStore := store.NewRedisStore(redisAddr)

	pol, err := policy.LoadPolicy("config/janus.json")

	if err != nil {
		log.Fatalf("Failed to load global config: %v", err)
	}
	log.Printf("Policy version: %d loaded", pol.Version)

	ac := admission.NewAdmissionController(pol, redisStore)


	// Init DB

	db.Init()
	defer db.Pool.Close()

	worker.StartJanusService(2, ac) // 2 worker threads only for dequeing from the queue. (async with Janus Singleton thread)
	worker.StartDBWriter(2)     // 2 worker thread to save the processed job into DB (async with Janus singleton thread)

	//Router
	mux := http.NewServeMux()

	// Route + middleware

	mux.Handle(
		"/dashboard/jobs",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobFromDashboard),
		),
	)
	mux.Handle(
		"/dashboard/jobs/batch",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobBatchFromDashboard), // partial
		),
	)

	mux.Handle(
		"/dashboard/jobs/batch/atomic",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobBatchAtomicFromDashboard), // atomic
		),
	)

	mux.Handle(
		"/system/jobs",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobFromSystem),
		),
	)

	mux.Handle(
		"/system/jobs/batch",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobBatchFromSystem), // partial
		),
	)

	mux.Handle(
		"/system/jobs/batch/atomic",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobBatchAtomicFromSystem), // atomic
		),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("HTTP server running on :8080")

	// Start server

	log.Println("Starting server...")

	err = server.ListenAndServe()

	log.Println("Server stopped") // runs only when server exits

	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

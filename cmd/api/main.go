package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/satyamraj1643/janus/db"
	"github.com/satyamraj1643/janus/handler"
	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/listener"
	"github.com/satyamraj1643/janus/middleware"
	"github.com/satyamraj1643/janus/worker"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Initialise admission controller
	// Initialise admission controller
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR not set")
	}

	redisStore := store.NewRedisStore(redisAddr)

	// Create a background context for initial ping
	if err := redisStore.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
	}

	log.Println("Connected to Redis")

	ac := admission.NewAdmissionController(redisStore)

	// Init DB

	db.Init()
	defer db.Pool.Close()

	dbURL := os.Getenv("DB_URL")
	listener.StartConfigListener(dbURL, redisStore)

	//worker.StartJanusService(2, ac) // Removed in favor of synchronous admission
	worker.StartDBWriter(2) // 2 worker thread to save the processed job into DB (async with Janus singleton thread)

	dashboardHandler := &handler.JobHandler{AC: ac, FromDashboard: true}
	systemHandler := &handler.JobHandler{AC: ac, FromDashboard: false}

	//Router
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Route + middleware

	mux.Handle(
		"/dashboard/jobs",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(dashboardHandler.CreateJob),
		),
	)
	mux.Handle(
		"/dashboard/jobs/batch",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(dashboardHandler.CreateJobBatch), // partial
		),
	)

	mux.Handle(
		"/dashboard/jobs/batch/atomic",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(dashboardHandler.CreateJobBatchAtomic), // atomic
		),
	)

	mux.Handle(
		"/system/jobs",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(systemHandler.CreateJob),
		),
	)

	mux.Handle(
		"/system/jobs/batch",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(systemHandler.CreateJobBatch), // partial
		),
	)

	mux.Handle(
		"/system/jobs/batch/atomic",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(systemHandler.CreateJobBatchAtomic), // atomic
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

	err := server.ListenAndServe()

	log.Println("Server stopped") // runs only when server exits

	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

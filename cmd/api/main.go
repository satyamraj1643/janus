package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/satyamraj1643/janus/db"
	"github.com/satyamraj1643/janus/handler"
	"github.com/satyamraj1643/janus/middleware"
	"github.com/satyamraj1643/janus/worker"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Init DB

	db.Init()
	defer db.Pool.Close()

	worker.Start(0) // 2 worker threads only for dequeing from the queue.

	//Router
	mux := http.NewServeMux()

	// Route + middleware

	mux.Handle(
		"POST /jobs/batch",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobBatch), // partial
		),
	)


	mux.Handle(
		"POST /jobs",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJob),
		),
	)

	
	mux.Handle(
		"POST /jobs/batch/atomic",
		middleware.ServiceRunningOnly(
			http.HandlerFunc(handler.CreateJobBatchAtomic), // atomic
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

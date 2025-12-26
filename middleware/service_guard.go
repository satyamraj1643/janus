package middleware

import (
	"log"
	"net/http"
	db "github.com/satyamraj1643/janus/db"
)

func ServiceRunningOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		log.Printf("userid: %s", userID)
		//userID := "2dad64a8-3f87-4d6e-9b4c-1cfa5917fd4b"

		if userID == "" {
			http.Error(w, "Missing user id", http.StatusBadRequest)
			return
		}

		running, err := db.IsServiceRunning(userID)

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if !running {
			http.Error(w, "Service is paused, please enable it from dashboard, then try running the service..", http.StatusForbidden)
			return
		}
        log.Printf("Forwarding request to respective handler ---")
		next.ServeHTTP(w,r)
	})
}
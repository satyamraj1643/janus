package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	db "github.com/satyamraj1643/janus/db"
	configStore "github.com/satyamraj1643/janus/globalStore"
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

		running, errService := db.IsServiceRunning(userID)

		if errService != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		if !running {
			http.Error(w, "Service is paused, please enable it from dashboard, then try running the service..", http.StatusForbidden)
			return
		}

		// 1. Try cache
		cached, found := configStore.Get(userID)

		var activeConfig json.RawMessage
		var activeConfigID string
		activeConfigExists := false

		if found {
			activeConfig = cached.Config
			activeConfigID = cached.ID
			activeConfigExists = true
		} else {
			// 2. Fetch from DB
			var errConfig error
			activeConfig, activeConfigID, activeConfigExists, errConfig = db.GetActiveJanusConfig(userID)
			if errConfig != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			// 3. Populate cache if found
			if activeConfigExists {
				configStore.Set(userID, activeConfig, activeConfigID)
			}
		}

		if !activeConfigExists {
			http.Error(w, "No active Janus config found, please create one from dashboard", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), activeConfigKey, activeConfig)
		ctx = context.WithValue(ctx, activeConfigIDKey, activeConfigID)
		ctx = context.WithValue(ctx, activeUserIDKey, userID)
		r = r.WithContext(ctx)

		log.Printf("Forwarding request to respective handler ---")
		next.ServeHTTP(w, r)
	})
}

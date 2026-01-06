package listener

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/satyamraj1643/janus/db"
	configStore "github.com/satyamraj1643/janus/globalStore"
	"github.com/satyamraj1643/janus/internal/store"
)

func StartConfigListener(dbURL string, s store.StateStore) {
	go func() {
		ctx := context.Background()

		conn, err := pgx.Connect(ctx, dbURL)
		if err != nil {
			log.Fatal("Failed to connect for LISTEN:", err)
		}

		_, err = conn.Exec(ctx, "LISTEN janus_config_update")
		if err != nil {
			log.Fatal("LISTEN failed:", err)
		}

		log.Println("Listening for janus_config_update")

		for {
			notification, err := conn.WaitForNotification(ctx)
			if err != nil {
				log.Println("Notification error:", err)
				continue
			}

			userID := notification.Payload
			log.Println("Config update for user:", userID)

			cfg, cfgID, ok, err := db.GetActiveJanusConfig(userID)
			if err != nil || !ok {
				configStore.Delete(userID)
				// If config is deleted/inactive, we definitely want to clear old quotas?
				// Maybe safe to flush here too, but critical when new config is loaded.
				continue
			}

			configStore.Set(userID, cfg, cfgID)

			// FLUSH Redis to ensure new limits take effect immediately
			log.Printf("Flushing Redis state for fresh start for user %s (actually global flush)", userID)
			if err := s.Flush(ctx); err != nil {
				log.Printf("Failed to flush redis: %v", err)
			}
		}
	}()
}

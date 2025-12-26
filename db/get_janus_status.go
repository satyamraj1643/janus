package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func IsServiceRunning(userID string) (bool, error) {
	var status string

	query := `
		SELECT status
		FROM service_status
		WHERE user_id = $1
	`

	err := Pool.QueryRow(
		context.Background(),
		query,
		userID,
	).Scan(&status)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("No service_status row for user:", userID)
			return false, nil // not running, but NOT an error
		}
		log.Println("DB error:", err)
		return false, err
	}

	log.Printf("Service status for user %s: %s", userID, status)
	return status == "running", nil
}


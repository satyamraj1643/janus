package db

import (
	"context"
	"encoding/json"
	"log"
)

func GetActiveJanusConfig(userID string) (json.RawMessage, string, bool, error) {
	var activeJanusConfig json.RawMessage
	var configID string

	query := `SELECT config, config_id from global_job_config where user_id = $1 and status = 'active'`

	err := Pool.QueryRow(
		context.Background(),
		query,
		userID,
	).Scan(&activeJanusConfig, &configID)

	if err != nil {
		log.Println("DB error:", err)
		return nil, "", false, err
	}

	return activeJanusConfig, configID, true, nil
}

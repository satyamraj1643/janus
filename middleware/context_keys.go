package middleware

import (
	"context"
	"encoding/json"
)

type contextKey string

const activeConfigKey contextKey = "activeJanusConfig"
const activeConfigIDKey contextKey = "activeJanusConfigID"
const activeUserIDKey contextKey = "activeUserID"

// GetActiveContext returns the active Janus config from the context
func GetActiveContext(ctx context.Context) (json.RawMessage, string, string, bool) {
	val := ctx.Value(activeConfigKey)
	idVal := ctx.Value(activeConfigIDKey)
	userVal := ctx.Value(activeUserIDKey)

	if val == nil {
		return nil, "", "", false
	}

	configID, _ := idVal.(string)
	userID, _ := userVal.(string)

	if cfg, ok := val.(json.RawMessage); ok {
		return cfg, configID, userID, true
	}

	if cfgPtr, ok := val.(*json.RawMessage); ok && cfgPtr != nil {
		return *cfgPtr, configID, userID, true
	}

	return nil, "", "", false
}

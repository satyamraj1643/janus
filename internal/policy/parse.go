package policy

import (
	"encoding/json"
	"fmt"
)

// ParseConfig unmarshals and validates a raw JSON config into a Policy struct
func ParseConfig(raw json.RawMessage) (*Policy, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty config")
	}

	var pol Policy
	if err := json.Unmarshal(raw, &pol); err != nil {
		return nil, fmt.Errorf("invalid config JSON: %w", err)
	}

	// Validation
	if pol.GlobalExecutionLimit.MaxJobs <= 0 {
		return nil, fmt.Errorf("max_jobs must be > 0")
	}
	if pol.GlobalExecutionLimit.WindowMs <= 0 {
		return nil, fmt.Errorf("window_ms must be > 0")
	}
	if pol.DefaultJobPolicy.IdempotencyWindowMs <= 0 {
		return nil, fmt.Errorf("idempotency_window_ms must be > 0")
	}

	return &pol, nil
}

package policy

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadPolicy(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	var p Policy
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse policy JSON: %w", err)
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}

	return &p, nil
}

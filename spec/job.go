package spec

import "time"

//Job represents a single execution request submitted to Janus
// It carries identity, classification, scope and business payload
// It does NOT carry execution or scheduling semantics

type Job struct {
	ID           string                 `json:"job_id"`
	TenantID     string                 `json:"tenant_id"`
	Priority     int                    `json:"priority"`
	Dependencies map[string]int         `json:"dependencies"`
	Payload      map[string]any `json:"payload"`


	// metadata (NOT user-provided)

	Source JobSource `json:"-"`
	BatchName string `json:"-"`
	BatchID string `json:"-"`
}

type JobDecision struct {
	// Identifiers
	JobID     string `json:"job_id"`
	BatchID   string `json:"batch_id"`
	BatchName string `json:"batch_name"`

	// Decision
	Status    string    `json:"status"` // accepted | rejected
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp"`

	// Full payload
	Job Job `json:"job"`
}

type JobSource string

const (
	JobSourceDashboard JobSource = "dashboard"
	JobSourceSystem JobSource = "system"
)

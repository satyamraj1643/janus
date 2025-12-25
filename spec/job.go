package spec

//Job represents a single execution request submitted to Janus
// It carries identity, classification, scope and business payload
// It does NOT carry execution or scheduling semantics

type Job struct {
	ID           string                 `json:"job_id"`
	TenantID     string                 `json:"tenant_id"`
	Priority     int                    `json:"priority"`
	Dependencies map[string]int         `json:"dependencies"`
	Payload      map[string]any `json:"payload"`
}

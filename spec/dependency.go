package spec

type DependencyType string

const (
	ExternalAPI     DependencyType = "external_api"
	InternalService DependencyType = "internal_service"
	Database        DependencyType = "database"
)

type RateLimit struct {
	MaxRequests int `json:"max_requests"`
	WindowMs    int `json:"window_ms"`
}

type Concurrency struct {
	MaxInflight int `json:"max_inflight"`
}
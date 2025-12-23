package spec

type OutcomeStatus string

const (
	OutcomeSuccess OutcomeStatus = "SUCCESS"
	OutcomeFailure OutcomeStatus = "FAILURE"
)

type ExecutionOutcome struct {
	JobID  string        `json:"job_id"`
	Status OutcomeStatus `json:"status"`
}

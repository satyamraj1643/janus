package handler

import "github.com/satyamraj1643/janus/spec"

type JobBatchRequest struct {
	BatchName string `json:"batch_name"`
	Jobs      []spec.Job `json:"jobs"`
}


type JobBatchResponse struct {
	BatchName string `json:"batch_name"`
	Status    string `json:"status"`   // full | partial | rejected
	Admitted  int    `json:"admitted"`
	Rejected  int    `json:"rejected"`
}


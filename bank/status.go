package bank

type status string

const (
	Pending   = "pending"
	Accepted  = "accepted"
	Processed = "processed"
	Canceled  = "canceled"
	Rejected  = "rejected"
	Errored   = "errored"
)

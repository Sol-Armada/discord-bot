package bank

type status string

const (
	Pending  = "pending"
	Accepted = "accepted"
	Rejected = "rejected"
	Errored  = "errored"
)

type Transaction struct {
	ID     string  `json:"id"`
	To     string  `json:"to"`
	From   string  `json:"from"`
	Amount float64 `json:"ammount"`
	Status status  `json:"status"`
}

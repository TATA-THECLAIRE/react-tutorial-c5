package models

type User struct {
}

type CreateTransactionPayload struct {
	Amount    float64 `json:"amount"`
	Type      string `json:"type"`
	Reason    string `json:"reason"`
}

type Transaction struct {
	ID        *int32 `json:"id"`
	Amount    float64 `json:"amount"`
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

const (
	TypeSaving     = "saving"
	TypeWithdrawal = "withdrawal"
)

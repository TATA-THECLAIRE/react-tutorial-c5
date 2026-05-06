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

type Account struct {
	AccountID string  `json:"account_id"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
}

type DashboardStats struct {
	Balance         float64 `json:"balance"`
	TotalSavings    float64 `json:"total_savings"`
	TotalWithdrawals float64 `json:"total_withdrawals"`
}

const (
	TypeSaving     = "saving"
	TypeWithdrawal = "withdrawal"
)

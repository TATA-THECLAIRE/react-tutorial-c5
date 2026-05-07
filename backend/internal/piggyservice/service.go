package piggyservice

import (
	"context"
	"fmt"

	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"piggy.com/internal/db/repo"
	"piggy.com/internal/db/sqlc"
	"piggy.com/internal/models"
)

type Service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) *Service {
	return &Service{repo: repo}
}

// define service methods here

func (s *Service) CreateTransaction(ctx context.Context, payload models.CreateTransactionPayload, userID pgtype.UUID) (*models.Transaction, error) {
		// float64 → pgtype.Numeric
	numeric := pgtype.Numeric{}
	numeric.Int = big.NewInt(int64(payload.Amount * 100))
	numeric.Exp = -2                   
	numeric.Valid = true

	// create the transaction
	
	transaction, err := s.repo.Do().CreateTransaction(ctx, sqlc.CreateTransactionParams{
		Amount: numeric,
		Type:   &payload.Type,
		Reason: &payload.Reason,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	// Update balance atomically
	if payload.Type == models.TypeSaving {
		_, err = s.repo.Do().IncrementBalance(ctx, sqlc.IncrementBalanceParams{
			Balance: numeric,
			UserID:  userID,
		})
	} else {
		// Check balance before decrementing
		account, err := s.repo.Do().GetAccountByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		balance := numericToFloat64(account.Balance)
		if payload.Amount > balance {
			return nil, fmt.Errorf("insufficient balance: have %.0f CFA, need %.0f CFA", balance, payload.Amount)
		}
		if balance-payload.Amount < 25 {
			return nil, fmt.Errorf("withdrawal would leave balance below 25 CFA")
		}
		_, err = s.repo.Do().DecrementBalance(ctx, sqlc.DecrementBalanceParams{
			Balance: numeric,
			UserID:  userID,
		})
		if err != nil {
			return nil, err
		}
	}

	return sqlCToAppTransaction(transaction), nil
}

func (s *Service) GetTransactions(ctx context.Context, userID pgtype.UUID, txnType string, size string) (*[]models.Transaction, error) {
	txns, err := s.repo.Do().GetTransactions(ctx,  userID)
	if err != nil {
		return nil, err
	}
	transactions := []models.Transaction{}
	for _, v := range txns {

		// filter by type if provided
		if txnType != "" && *v.Type != txnType {
			continue
		}
		transactions = append(transactions, *sqlCToAppTransaction(v))
		

	}
	// limit by size if provided
	if size != "" {
		n := 0
		fmt.Sscanf(size, "%d", &n)
		if n > 0 && n < len(transactions) {
			transactions = transactions[:n]
		}
	}
	return &transactions, nil
}

// GetDashboardStats returns balance + totals efficiently
func (s *Service) GetDashboardStats(ctx context.Context, userID pgtype.UUID) (*models.DashboardStats, error) {
	account, err := s.repo.Do().GetAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get totals by type — only two aggregated rows, not all transactions
	txns, err := s.repo.Do().GetTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalSavings, totalWithdrawals float64
	for _, t := range txns {
		amt := numericToFloat64(t.Amount)
		if t.Type != nil && *t.Type == models.TypeSaving {
			totalSavings += amt
		} else {
			totalWithdrawals += amt
		}
	}

	return &models.DashboardStats{
		Balance:          numericToFloat64(account.Balance),
		TotalSavings:     totalSavings,
		TotalWithdrawals: totalWithdrawals,
	}, nil
}

// Also create account when user registers
func (s *Service) CreateAccount(ctx context.Context, userID pgtype.UUID) error {
	_, err := s.repo.Do().CreateAccount(ctx, userID)
	return err
}

// helpers
func float64ToNumeric(f float64) pgtype.Numeric {
	bf := new(big.Float).SetFloat64(f)
	bf.Mul(bf, new(big.Float).SetFloat64(100))
	i, _ := bf.Int(nil)
	return pgtype.Numeric{Int: i, Exp: -2, Valid: true}
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	f, _ := new(big.Float).SetInt(n.Int).Float64()
	exp := n.Exp
	if exp < 0 {
		divisor := new(big.Float).SetFloat64(math.Pow(10, float64(-exp)))
		f, _ = new(big.Float).Quo(new(big.Float).SetFloat64(f), divisor).Float64()
	}
	return f
}

func sqlCToAppTransaction(t sqlc.Transaction) *models.Transaction {
	var amount float64

	if t.Amount.Valid && t.Amount.Int != nil {
		f, _ := new(big.Float).SetInt(t.Amount.Int).Float64()
		exp := t.Amount.Exp
		if exp < 0 {
			divisor := new(big.Float).SetFloat64(math.Pow(10, float64(-exp)))
			f, _ = new(big.Float).Quo(new(big.Float).SetFloat64(f), divisor).Float64()
		}
		amount = f
	}


	
	return &models.Transaction{
		Amount:    amount,
		Reason:    *t.Reason,
		Type:      *t.Type,
		ID:        &t.ID,
		CreatedAt: t.CreatedAt.Time.String(),
	}
}

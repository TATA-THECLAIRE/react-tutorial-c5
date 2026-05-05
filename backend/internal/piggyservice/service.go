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

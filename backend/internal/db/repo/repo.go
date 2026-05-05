package repo

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pgx5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"piggy.com/internal/db/sqlc"
	
)

type Repository interface {
	Begin(ctx context.Context) (sqlc.Querier, pgx5.Tx, error)
	Do() sqlc.Querier
}

type Impl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Impl {
	return &Impl{db: db}
}

func (u *Impl) Begin(ctx context.Context) (sqlc.Querier, pgx5.Tx, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return nil, tx, err
	}
	return sqlc.New(tx), tx, nil
}

func (u *Impl) Do() sqlc.Querier {
	return sqlc.New(u.db)
}

// MigrateUp function applies migrations to the database.
func MigrateUp(dbURL string, migrationsPath string, logger zerolog.Logger) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	// Create a new migration instance with the absolute path
	sourceURL := "file://" + filepath.ToSlash(absPath)
	m, err := migrate.New(sourceURL,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("error occured: %w", err)
	}
	defer m.Close()

	// Apply migrations
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// MigrateDown function rolls back migrations from the database.
func MigrateDown(dbURL string, migrationsPath string, logger zerolog.Logger) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	// Create a new migration instance with the absolute path
	m, err := migrate.New(
		"file://"+absPath,
		dbURL,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	// Apply migrations
	if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

package piggyservice

import (
    "context"
    "piggy.com/internal/db/sqlc"
)

func (s *Service) CreateUser(ctx context.Context, username, email, hashedPassword string) (sqlc.User, error) {
    return s.repo.Do().CreateUser(ctx, sqlc.CreateUserParams{
        Username: username,
        Email:    email,
        Password: hashedPassword,
    })
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (sqlc.User, error) {
    return s.repo.Do().GetUserByUsername(ctx, username)
}
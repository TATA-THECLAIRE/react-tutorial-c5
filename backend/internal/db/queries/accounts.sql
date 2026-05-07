-- name: CreateAccount :one
INSERT INTO user_accounts (user_id)
VALUES ($1)
RETURNING *;

-- name: GetAccountByUserID :one
SELECT * FROM user_accounts
WHERE user_id = $1;

-- name: IncrementBalance :one
UPDATE user_accounts
SET balance = balance + $1, updated_at = NOW()
WHERE user_id = $2
RETURNING *;

-- name: DecrementBalance :one
UPDATE user_accounts
SET balance = balance - $1, updated_at = NOW()
WHERE user_id = $2
RETURNING *;
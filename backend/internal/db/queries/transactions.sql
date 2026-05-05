-- name: CreateTransaction :one
INSERT INTO transactions (amount, reason, type, user_id) 
VALUES ($1, $2, $3, $4)
 RETURNING *;

-- name: GetTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC;
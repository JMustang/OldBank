-- name: CreateTransfer :one
INSERT INTO transfers (
    fromAccountId,
    toAccountId,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE
    fromAccountId = $1 OR
    toAccountId = $2
ORDER BY id
LIMIT $3
OFFSET $4;
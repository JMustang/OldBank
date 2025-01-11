-- name: CreateEntry :one
INSERT INTO entries (
    accoubt_id,
    amount
) VALUES (
    $1, $2
) RETURNING *;
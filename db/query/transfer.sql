-- name: CreateTransfer :one
INSERT INTO transfers (
    from_accoubt_id,
    to_accoubt_id,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;
-- name: CreateUser :one
INSERT INTO "user" (
  email, password
) VALUES (
  $1, $2
)
RETURNING *;
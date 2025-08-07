-- name: CreateUser :one
INSERT INTO
  "user" (email, hash)
VALUES
  ($1, $2)
RETURNING
  id;

-- name: GetUsers :many
SELECT
  u.*
FROM
  "user" u;

-- name: GetUserByEmail :one
SELECT
  u.*
FROM
  "user" u
WHERE
  u.email = $1;

-- name: GetUserById :one
SELECT
  u.*
FROM
  "user" u
WHERE
  u.id = $1;

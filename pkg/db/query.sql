-- name: CreateUser :one
INSERT INTO "user" (email, hash) 
VALUES ($1, $2)
RETURNING id;

-- name: GetUserIdByCredentials :one
SELECT u.id
FROM "user" u
WHERE u.email = $1
AND u.hash = $2;

-- name: GetUserById :one
SELECT *
FROM "user" u
WHERE u.id = $1;
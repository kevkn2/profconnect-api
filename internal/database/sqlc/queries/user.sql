-- name: GetUserByID :one
SELECT id, name, email, hashed_password, role FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, hashed_password, role FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, hashed_password, role) 
VALUES ($1, $2, $3, $4) 
RETURNING id, name, email, hashed_password, role;

-- name: UpdateUser :one
UPDATE users 
SET name = $1, email = $2, hashed_password = $3, role = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $5
RETURNING id, name, email, hashed_password, role;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

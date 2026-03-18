-- name: GetUserByEmail :one
SELECT id, name, email, hashed_password FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, hashed_password) 
VALUES ($1, $2, $3) 
RETURNING id, name, email, hashed_password;

-- name: UpdateUser :one
UPDATE users 
SET name = $1, email = $2, hashed_password = $3 
WHERE id = $4 
RETURNING id, name, email, hashed_password;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

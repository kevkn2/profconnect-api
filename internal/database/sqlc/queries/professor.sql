-- name: CreateProfessor :one
INSERT INTO professors (user_id, university, department) 
VALUES ($1, $2, $3) 
RETURNING id, user_id, university, department;

-- name: GetProfessorByUserID :one
SELECT id, user_id, university, department FROM professors WHERE user_id = $1;

-- name: UpdateProfessor :one
UPDATE professors
SET university = $1, department = $2, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $3
RETURNING id, user_id, university, department;
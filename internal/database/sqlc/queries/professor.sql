-- name: CreateProfessor :one
INSERT INTO professors (user_id, university, department) 
VALUES ($1, $2, $3) 
RETURNING id, user_id, university, department;

-- name: GetProfessorByUserID :one
SELECT 
    p.id,
    p.user_id,
    u.name AS user_name,
    u.email AS user_email,
    u.role AS user_role,
    p.university,
    p.department 
FROM professors p
JOIN users u ON p.user_id = u.id
WHERE user_id = $1 AND NOT u.deleted;

-- name: UpdateProfessor :one
WITH updated AS (
    UPDATE professors
    SET university = $1,
        department = $2,
        updated_at = CURRENT_TIMESTAMP
    WHERE user_id = $3
    RETURNING id, user_id, university, department
)
SELECT 
    u.id,
    u.user_id,
    u.university,
    u.department,
    usr.name,
    usr.email,
    usr.role
FROM updated u
JOIN users usr ON usr.id = u.user_id;

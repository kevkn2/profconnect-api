-- name: CreateStudent :one
INSERT INTO students (user_id, university, department, research_interests) 
VALUES ($1, $2, $3, $4) 
RETURNING id, user_id, university, department, research_interests;

-- name: GetStudentByUserID :one
SELECT
    s.id,
    s.user_id,
    u.name AS user_name,
    u.email AS user_email,
    u.role AS user_role,
    s.university,
    s.department,
    s.research_interests
FROM students s
JOIN users u ON s.user_id = u.id
WHERE user_id = $1 AND NOT u.deleted;

-- name: GetStudentByID :one
SELECT
    s.id,
    s.user_id,
    u.name AS user_name,
    u.email AS user_email,
    u.role AS user_role,
    s.university,
    s.department,
    s.research_interests
FROM students s
JOIN users u ON s.user_id = u.id
WHERE s.id = $1 AND NOT u.deleted;

-- name: ListStudents :many
SELECT
    s.id,
    s.user_id,
    u.name AS user_name,
    u.email AS user_email,
    s.university,
    s.department,
    s.research_interests
FROM students s
JOIN users u ON s.user_id = u.id
WHERE NOT u.deleted
ORDER BY u.name ASC;

-- name: UpdateStudent :one
WITH updated AS (
    UPDATE students
    SET university = $1,
        department = $2,
        research_interests = $3,
        updated_at = CURRENT_TIMESTAMP
    WHERE user_id = $4
    RETURNING id, user_id, university, department, research_interests
)
SELECT 
    u.id,
    u.user_id,
    u.university,
    u.department,
    u.research_interests,
    usr.name,
    usr.email,
    usr.role
FROM updated u
JOIN users usr ON usr.id = u.user_id;

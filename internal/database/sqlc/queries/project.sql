-- name: CreateProject :one
INSERT INTO projects (professor_id, title, description, slots, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, professor_id, title, description, slots, status;

-- name: GetProjectByID :one
SELECT
    p.id,
    p.professor_id,
    p.title,
    p.description,
    p.slots,
    p.status,
    prof.user_id AS professor_user_id,
    u.name AS professor_name,
    u.email AS professor_email,
    prof.university AS professor_university,
    prof.department AS professor_department
FROM projects p
JOIN professors prof ON p.professor_id = prof.id
JOIN users u ON prof.user_id = u.id
WHERE p.id = $1;

-- name: ListProjects :many
SELECT
    p.id,
    p.professor_id,
    p.title,
    p.description,
    p.slots,
    p.status,
    prof.user_id AS professor_user_id,
    u.name AS professor_name,
    u.email AS professor_email,
    prof.university AS professor_university,
    prof.department AS professor_department
FROM projects p
JOIN professors prof ON p.professor_id = prof.id
JOIN users u ON prof.user_id = u.id
ORDER BY p.created_at DESC;

-- name: ListProjectsByProfessor :many
SELECT
    p.id,
    p.professor_id,
    p.title,
    p.description,
    p.slots,
    p.status,
    prof.user_id AS professor_user_id,
    u.name AS professor_name,
    u.email AS professor_email,
    prof.university AS professor_university,
    prof.department AS professor_department
FROM projects p
JOIN professors prof ON p.professor_id = prof.id
JOIN users u ON prof.user_id = u.id
WHERE p.professor_id = $1
ORDER BY p.created_at DESC;

-- name: UpdateProjectStatus :exec
UPDATE projects
SET status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: CountApprovedApplications :one
SELECT COUNT(*) AS approved_count
FROM project_applications
WHERE project_id = $1 AND status = 'approved';

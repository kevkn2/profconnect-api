-- name: CreateProjectApplication :one
INSERT INTO project_applications (project_id, student_id, status, message)
VALUES ($1, $2, $3, $4)
RETURNING id, project_id, student_id, status, message;

-- name: GetProjectApplicationByID :one
SELECT
    pa.id,
    pa.project_id,
    pa.student_id,
    pa.status,
    pa.message,
    s.user_id AS student_user_id,
    u.name AS student_name,
    u.email AS student_email,
    s.university AS student_university,
    s.department AS student_department,
    s.research_interests AS student_research_interests
FROM project_applications pa
JOIN students s ON pa.student_id = s.id
JOIN users u ON s.user_id = u.id
WHERE pa.id = $1;

-- name: GetProjectApplicationByProjectAndStudent :one
SELECT id, project_id, student_id, status, message
FROM project_applications
WHERE project_id = $1 AND student_id = $2;

-- name: ListApplicationsByProject :many
SELECT
    pa.id,
    pa.project_id,
    pa.student_id,
    pa.status,
    pa.message,
    s.user_id AS student_user_id,
    u.name AS student_name,
    u.email AS student_email,
    s.university AS student_university,
    s.department AS student_department,
    s.research_interests AS student_research_interests
FROM project_applications pa
JOIN students s ON pa.student_id = s.id
JOIN users u ON s.user_id = u.id
WHERE pa.project_id = $1
ORDER BY pa.created_at DESC;

-- name: ListApplicationsByStudent :many
SELECT
    pa.id,
    pa.project_id,
    pa.student_id,
    pa.status,
    pa.message,
    p.title AS project_title,
    p.description AS project_description,
    p.status AS project_status
FROM project_applications pa
JOIN projects p ON pa.project_id = p.id
WHERE pa.student_id = $1
ORDER BY pa.created_at DESC;

-- name: CheckApplicationsByProjectAndStudent :one
SELECT EXISTS (
    SELECT 1 FROM project_applications
    WHERE student_id = $1
        AND project_id = $2
        AND status IN ('pending', 'approved')
);

-- name: UpdateProjectApplicationStatus :one
UPDATE project_applications
SET status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING id, project_id, student_id, status, message;

-- name: DeleteProjectApplication :exec
DELETE FROM project_applications
WHERE id = $1;

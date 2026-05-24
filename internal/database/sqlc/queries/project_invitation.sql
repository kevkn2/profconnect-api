-- name: CreateProjectInvitation :one
INSERT INTO project_invitations (project_id, student_id, status, message)
VALUES ($1, $2, $3, $4)
RETURNING id, project_id, student_id, status, message, responded_at;

-- name: GetProjectInvitationByID :one
SELECT
    pi.id,
    pi.project_id,
    pi.student_id,
    pi.status,
    pi.message,
    pi.responded_at,
    s.user_id AS student_user_id,
    u.name AS student_name,
    u.email AS student_email,
    s.university AS student_university,
    s.department AS student_department,
    s.research_interests AS student_research_interests
FROM project_invitations pi
JOIN students s ON pi.student_id = s.id
JOIN users u ON s.user_id = u.id
WHERE pi.id = $1;

-- name: GetInvitationByProjectAndStudent :one
SELECT id, project_id, student_id, status, message, responded_at
FROM project_invitations
WHERE project_id = $1 AND student_id = $2;

-- name: ListInvitationsByProject :many
SELECT
    pi.id,
    pi.project_id,
    pi.student_id,
    pi.status,
    pi.message,
    pi.responded_at,
    s.user_id AS student_user_id,
    u.name AS student_name,
    u.email AS student_email,
    s.university AS student_university,
    s.department AS student_department,
    s.research_interests AS student_research_interests
FROM project_invitations pi
JOIN students s ON pi.student_id = s.id
JOIN users u ON s.user_id = u.id
WHERE pi.project_id = $1
ORDER BY pi.created_at DESC;

-- name: ListInvitationsByStudent :many
SELECT
    pi.id,
    pi.project_id,
    pi.student_id,
    pi.status,
    pi.message,
    pi.responded_at,
    p.title AS project_title,
    p.description AS project_description,
    p.status AS project_status
FROM project_invitations pi
JOIN projects p ON pi.project_id = p.id
WHERE pi.student_id = $1
ORDER BY pi.created_at DESC;

-- name: UpdateProjectInvitationStatus :one
UPDATE project_invitations
SET status = $1,
    responded_at = CASE WHEN $1 IN ('accepted', 'declined') THEN CURRENT_TIMESTAMP ELSE responded_at END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING id, project_id, student_id, status, message, responded_at;

-- name: DeleteProjectInvitation :exec
DELETE FROM project_invitations
WHERE id = $1;

-- name: CreateProjectMember :one
INSERT INTO project_members (project_id, student_id, source, source_ref_id, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, student_id, source, source_ref_id, status, joined_at, left_at;

-- name: GetProjectMemberByID :one
SELECT
    pm.id,
    pm.project_id,
    pm.student_id,
    pm.source,
    pm.source_ref_id,
    pm.status,
    pm.joined_at,
    pm.left_at,
    s.user_id AS student_user_id,
    u.name AS student_name,
    u.email AS student_email,
    s.university AS student_university,
    s.department AS student_department,
    s.research_interests AS student_research_interests
FROM project_members pm
JOIN students s ON pm.student_id = s.id
JOIN users u ON s.user_id = u.id
WHERE pm.id = $1;

-- name: GetActiveMemberByProjectAndStudent :one
SELECT id, project_id, student_id, source, source_ref_id, status, joined_at, left_at
FROM project_members
WHERE project_id = $1 AND student_id = $2 AND status = 'active';

-- name: ListMembersByProject :many
SELECT
    pm.id,
    pm.project_id,
    pm.student_id,
    pm.source,
    pm.source_ref_id,
    pm.status,
    pm.joined_at,
    pm.left_at,
    s.user_id AS student_user_id,
    u.name AS student_name,
    u.email AS student_email,
    s.university AS student_university,
    s.department AS student_department,
    s.research_interests AS student_research_interests
FROM project_members pm
JOIN students s ON pm.student_id = s.id
JOIN users u ON s.user_id = u.id
WHERE pm.project_id = $1
ORDER BY pm.joined_at DESC;

-- name: ListActiveMembershipsByStudent :many
SELECT
    pm.id,
    pm.project_id,
    pm.student_id,
    pm.source,
    pm.source_ref_id,
    pm.status,
    pm.joined_at,
    pm.left_at,
    p.title AS project_title,
    p.description AS project_description,
    p.status AS project_status,
    p.slots AS project_slots
FROM project_members pm
JOIN projects p ON pm.project_id = p.id
WHERE pm.student_id = $1 AND pm.status = 'active'
ORDER BY pm.joined_at DESC;

-- name: CountActiveMembersByProject :one
SELECT COUNT(*) AS active_count
FROM project_members
WHERE project_id = $1 AND status = 'active';

-- name: UpdateProjectMemberStatus :one
UPDATE project_members
SET status = $1,
    left_at = CASE WHEN $1 IN ('removed', 'left') THEN CURRENT_TIMESTAMP ELSE left_at END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING id, project_id, student_id, source, source_ref_id, status, joined_at, left_at;

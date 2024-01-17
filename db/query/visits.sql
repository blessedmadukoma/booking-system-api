-- name: CreateVisit :one
INSERT INTO visits (
 status,
 reason,
 employee_id,
 visitor_id
) VALUES (
 $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateVisit :one
UPDATE visits
SET
status = $4,
reason = $5,
updated_at = now()
WHERE id = $1
AND visitor_id = $2
AND employee_id = $3
RETURNING *;

-- name: GetVisitsByID :one
SELECT * FROM visits 
WHERE id = $1 LIMIT 1;

-- name: ListVisits :many
SELECT v.id, v.status, v.reason, v.employee_id, v.visitor_id, v.created_at, v.updated_at, vi.fullname as visitor_name, vi.email as visitor_email, vi.mobile as visitor_mobile, vi.company_name, vi.picture as visitor_picture, vi.sign_in, vi.sign_out, e.fullname as employee_name, e.email as employee_email
FROM visits v
LEFT JOIN visitors vi ON vi.id = v.visitor_id
LEFT JOIN employees e ON e.id = v.employee_id
ORDER BY v.created_at DESC;

-- name: ListVisitsByStatus :many
SELECT * FROM visits 
WHERE status = $1
ORDER BY id;

-- name: GetVisitByVisitorID :one
SELECT * FROM visits 
WHERE visitor_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: ListEmployeeVisits :many
SELECT visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visits.status, visits.reason, visits.employee_id, visits.visitor_id, visits.created_at, visits.updated_at, visits.id
FROM visits
INNER JOIN visitors ON visits.visitor_id=visitors.id
WHERE visits.employee_id=$1
ORDER BY visits.created_at DESC;

-- name: DeleteVisit :exec
DELETE FROM visits
WHERE id = $1;
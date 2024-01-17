-- name: CreateVisitor :one
INSERT INTO visitors (
 fullname,
 email,
 mobile,
 company_name,
 picture,
 employee_id
) VALUES (
 $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetVisitorByID :one
SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id, employees.fullname AS employee_name, employees.email AS employee_email
FROM visitors
LEFT JOIN employees ON visitors.employee_id=employees.id
WHERE visitors.id = $1 LIMIT 1;
-- SELECT * FROM visitors 
-- WHERE id = $1 LIMIT 1;

-- -- name: ListVisitors :many
-- SELECT * FROM visitors 
-- ORDER BY id
-- LIMIT $1
-- OFFSET $2;

-- name: ListVisitors :many
SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id, employees.fullname AS employee_name, employees.email AS employee_email
FROM visitors
LEFT JOIN employees ON visitors.employee_id=employees.id
WHERE sign_in >= NOW() - INTERVAL '24 HOURS'
AND sign_out IS NULL
ORDER by visitors.id;

-- name: ListAllVisitors :many
SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id, employees.fullname AS employee_name, employees.email AS employee_email
FROM visitors
LEFT JOIN employees ON visitors.employee_id=employees.id
-- ORDER by visitors.id;
ORDER by visitors.created_at DESC;

-- name: UpdateVisitor :one
UPDATE visitors
SET
sign_out = now(),
updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteVisitor :exec
DELETE FROM visitors
WHERE id = $1;
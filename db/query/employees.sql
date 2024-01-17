-- name: CreateEmployee :one
INSERT INTO employees (
 fullname, email, mobile, password, token
) VALUES (
 $1, $2, $3, $4, $5
)
RETURNING *;

-- name: LoginEmployee :one
UPDATE employees
SET token = $2,
updated_at = now()
WHERE email = $1
RETURNING *;

-- name: GetEmployeeByID :one
SELECT * FROM employees 
WHERE id = $1 LIMIT 1;

-- name: GetEmployeeByEmail :one
SELECT * FROM employees 
WHERE email = $1 LIMIT 1;

-- name: ListEmployees :many
SELECT * FROM employees 
ORDER BY id;

-- -- name: ListEmployees :many
-- SELECT * FROM employees 
-- ORDER BY id
-- LIMIT $1
-- OFFSET $2;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;

-- name: GetEmployeeVisitors :many
SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id
FROM employees
LEFT JOIN visitors ON employees.id=visitors.employee_id
WHERE employees.id = $1
AND sign_in >= NOW() - INTERVAL '24 HOURS';
-- SELECT employees.id, employees.fullname AS employee_name, employees.email AS employee_email, visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id
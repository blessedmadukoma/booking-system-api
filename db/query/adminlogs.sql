-- name: CreateAdminLog :one
INSERT INTO adminlogs (
 admin_id
) VALUES (
 $1
)
RETURNING *;

-- name: UpdateAdminLogs :one
UPDATE adminlogs
SET logged_out = now()
WHERE admin_id = $1
RETURNING *;

-- name: GetAdminLogsByAdminID :one
SELECT * FROM adminlogs
WHERE admin_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAdminLogByAdminID :many
SELECT * FROM adminlogs
WHERE admin_id = $1
ORDER BY id;

-- name: ListAllAdminLogs :many
SELECT * FROM adminlogs
ORDER BY id;
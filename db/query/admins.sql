-- name: CreateAdmin :one
INSERT INTO admin (
 fullname, email, role, password
) VALUES (
 $1, $2, $3, $4
)
RETURNING *;

-- name: LoginAdmin :one
UPDATE admin
SET logged_in = now()
WHERE email = $1
RETURNING *;

-- name: LogoutAdmin :one
UPDATE admin
SET logged_out = now()
WHERE email = $1
RETURNING *;

-- name: GetAdminByID :one
SELECT * FROM admin 
WHERE id = $1 LIMIT 1;

-- name: GetAdminByEmail :one
SELECT * FROM admin 
WHERE email = $1 LIMIT 1;

-- name: ListAdmins :many
SELECT * FROM admin 
ORDER BY id;

-- name: DeleteAdmin :exec
DELETE FROM admin
WHERE id = $1;
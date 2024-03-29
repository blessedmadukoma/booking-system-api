// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: admins.sql

package db

import (
	"context"
	"database/sql"
)

const createAdmin = `-- name: CreateAdmin :one
INSERT INTO admin (
 fullname, email, role, password
) VALUES (
 $1, $2, $3, $4
)
RETURNING id, fullname, email, role, logged_in, logged_out, password, created_at, updated_at
`

type CreateAdminParams struct {
	Fullname string         `json:"fullname"`
	Email    string         `json:"email"`
	Role     sql.NullString `json:"role"`
	Password string         `json:"password"`
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) (Admin, error) {
	row := q.db.QueryRowContext(ctx, createAdmin,
		arg.Fullname,
		arg.Email,
		arg.Role,
		arg.Password,
	)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Role,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAdmin = `-- name: DeleteAdmin :exec
DELETE FROM admin
WHERE id = $1
`

func (q *Queries) DeleteAdmin(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAdmin, id)
	return err
}

const getAdminByEmail = `-- name: GetAdminByEmail :one
SELECT id, fullname, email, role, logged_in, logged_out, password, created_at, updated_at FROM admin 
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetAdminByEmail(ctx context.Context, email string) (Admin, error) {
	row := q.db.QueryRowContext(ctx, getAdminByEmail, email)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Role,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAdminByID = `-- name: GetAdminByID :one
SELECT id, fullname, email, role, logged_in, logged_out, password, created_at, updated_at FROM admin 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAdminByID(ctx context.Context, id int64) (Admin, error) {
	row := q.db.QueryRowContext(ctx, getAdminByID, id)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Role,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAdmins = `-- name: ListAdmins :many
SELECT id, fullname, email, role, logged_in, logged_out, password, created_at, updated_at FROM admin 
ORDER BY id
`

func (q *Queries) ListAdmins(ctx context.Context) ([]Admin, error) {
	rows, err := q.db.QueryContext(ctx, listAdmins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Admin{}
	for rows.Next() {
		var i Admin
		if err := rows.Scan(
			&i.ID,
			&i.Fullname,
			&i.Email,
			&i.Role,
			&i.LoggedIn,
			&i.LoggedOut,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const loginAdmin = `-- name: LoginAdmin :one
UPDATE admin
SET logged_in = now()
WHERE email = $1
RETURNING id, fullname, email, role, logged_in, logged_out, password, created_at, updated_at
`

func (q *Queries) LoginAdmin(ctx context.Context, email string) (Admin, error) {
	row := q.db.QueryRowContext(ctx, loginAdmin, email)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Role,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const logoutAdmin = `-- name: LogoutAdmin :one
UPDATE admin
SET logged_out = now()
WHERE email = $1
RETURNING id, fullname, email, role, logged_in, logged_out, password, created_at, updated_at
`

func (q *Queries) LogoutAdmin(ctx context.Context, email string) (Admin, error) {
	row := q.db.QueryRowContext(ctx, logoutAdmin, email)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Role,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

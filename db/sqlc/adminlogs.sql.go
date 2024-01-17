// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: adminlogs.sql

package db

import (
	"context"
)

const createAdminLog = `-- name: CreateAdminLog :one
INSERT INTO adminlogs (
 admin_id
) VALUES (
 $1
)
RETURNING id, admin_id, logged_in, logged_out, created_at
`

func (q *Queries) CreateAdminLog(ctx context.Context, adminID int64) (Adminlog, error) {
	row := q.db.QueryRowContext(ctx, createAdminLog, adminID)
	var i Adminlog
	err := row.Scan(
		&i.ID,
		&i.AdminID,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.CreatedAt,
	)
	return i, err
}

const getAdminLogByAdminID = `-- name: GetAdminLogByAdminID :many
SELECT id, admin_id, logged_in, logged_out, created_at FROM adminlogs
WHERE admin_id = $1
ORDER BY id
`

func (q *Queries) GetAdminLogByAdminID(ctx context.Context, adminID int64) ([]Adminlog, error) {
	rows, err := q.db.QueryContext(ctx, getAdminLogByAdminID, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Adminlog{}
	for rows.Next() {
		var i Adminlog
		if err := rows.Scan(
			&i.ID,
			&i.AdminID,
			&i.LoggedIn,
			&i.LoggedOut,
			&i.CreatedAt,
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

const getAdminLogsByAdminID = `-- name: GetAdminLogsByAdminID :one
SELECT id, admin_id, logged_in, logged_out, created_at FROM adminlogs
WHERE admin_id = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetAdminLogsByAdminID(ctx context.Context, adminID int64) (Adminlog, error) {
	row := q.db.QueryRowContext(ctx, getAdminLogsByAdminID, adminID)
	var i Adminlog
	err := row.Scan(
		&i.ID,
		&i.AdminID,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.CreatedAt,
	)
	return i, err
}

const listAllAdminLogs = `-- name: ListAllAdminLogs :many
SELECT id, admin_id, logged_in, logged_out, created_at FROM adminlogs
ORDER BY id
`

func (q *Queries) ListAllAdminLogs(ctx context.Context) ([]Adminlog, error) {
	rows, err := q.db.QueryContext(ctx, listAllAdminLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Adminlog{}
	for rows.Next() {
		var i Adminlog
		if err := rows.Scan(
			&i.ID,
			&i.AdminID,
			&i.LoggedIn,
			&i.LoggedOut,
			&i.CreatedAt,
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

const updateAdminLogs = `-- name: UpdateAdminLogs :one
UPDATE adminlogs
SET logged_out = now()
WHERE admin_id = $1
RETURNING id, admin_id, logged_in, logged_out, created_at
`

func (q *Queries) UpdateAdminLogs(ctx context.Context, adminID int64) (Adminlog, error) {
	row := q.db.QueryRowContext(ctx, updateAdminLogs, adminID)
	var i Adminlog
	err := row.Scan(
		&i.ID,
		&i.AdminID,
		&i.LoggedIn,
		&i.LoggedOut,
		&i.CreatedAt,
	)
	return i, err
}

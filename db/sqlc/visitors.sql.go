// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: visitors.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createVisitor = `-- name: CreateVisitor :one
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
RETURNING id, fullname, email, mobile, company_name, picture, sign_in, sign_out, employee_id, created_at, updated_at
`

type CreateVisitorParams struct {
	Fullname    string        `json:"fullname"`
	Email       string        `json:"email"`
	Mobile      string        `json:"mobile"`
	CompanyName string        `json:"company_name"`
	Picture     string        `json:"picture"`
	EmployeeID  sql.NullInt64 `json:"employee_id"`
}

func (q *Queries) CreateVisitor(ctx context.Context, arg CreateVisitorParams) (Visitor, error) {
	row := q.db.QueryRowContext(ctx, createVisitor,
		arg.Fullname,
		arg.Email,
		arg.Mobile,
		arg.CompanyName,
		arg.Picture,
		arg.EmployeeID,
	)
	var i Visitor
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Mobile,
		&i.CompanyName,
		&i.Picture,
		&i.SignIn,
		&i.SignOut,
		&i.EmployeeID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteVisitor = `-- name: DeleteVisitor :exec
DELETE FROM visitors
WHERE id = $1
`

func (q *Queries) DeleteVisitor(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteVisitor, id)
	return err
}

const getVisitorByID = `-- name: GetVisitorByID :one
SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id, employees.fullname AS employee_name, employees.email AS employee_email
FROM visitors
LEFT JOIN employees ON visitors.employee_id=employees.id
WHERE visitors.id = $1 LIMIT 1
`

type GetVisitorByIDRow struct {
	ID            int64          `json:"id"`
	Fullname      string         `json:"fullname"`
	Email         string         `json:"email"`
	Mobile        string         `json:"mobile"`
	CompanyName   string         `json:"company_name"`
	Picture       string         `json:"picture"`
	SignIn        time.Time      `json:"sign_in"`
	SignOut       sql.NullTime   `json:"sign_out"`
	EmployeeID    sql.NullInt64  `json:"employee_id"`
	EmployeeName  sql.NullString `json:"employee_name"`
	EmployeeEmail sql.NullString `json:"employee_email"`
}

func (q *Queries) GetVisitorByID(ctx context.Context, id int64) (GetVisitorByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getVisitorByID, id)
	var i GetVisitorByIDRow
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Mobile,
		&i.CompanyName,
		&i.Picture,
		&i.SignIn,
		&i.SignOut,
		&i.EmployeeID,
		&i.EmployeeName,
		&i.EmployeeEmail,
	)
	return i, err
}

const listAllVisitors = `-- name: ListAllVisitors :many
SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id, employees.fullname AS employee_name, employees.email AS employee_email
FROM visitors
LEFT JOIN employees ON visitors.employee_id=employees.id
ORDER by visitors.created_at DESC
`

type ListAllVisitorsRow struct {
	ID            int64          `json:"id"`
	Fullname      string         `json:"fullname"`
	Email         string         `json:"email"`
	Mobile        string         `json:"mobile"`
	CompanyName   string         `json:"company_name"`
	Picture       string         `json:"picture"`
	SignIn        time.Time      `json:"sign_in"`
	SignOut       sql.NullTime   `json:"sign_out"`
	EmployeeID    sql.NullInt64  `json:"employee_id"`
	EmployeeName  sql.NullString `json:"employee_name"`
	EmployeeEmail sql.NullString `json:"employee_email"`
}

// ORDER by visitors.id;
func (q *Queries) ListAllVisitors(ctx context.Context) ([]ListAllVisitorsRow, error) {
	rows, err := q.db.QueryContext(ctx, listAllVisitors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAllVisitorsRow{}
	for rows.Next() {
		var i ListAllVisitorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Fullname,
			&i.Email,
			&i.Mobile,
			&i.CompanyName,
			&i.Picture,
			&i.SignIn,
			&i.SignOut,
			&i.EmployeeID,
			&i.EmployeeName,
			&i.EmployeeEmail,
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

const listVisitors = `-- name: ListVisitors :many


SELECT visitors.id, visitors.fullname, visitors.email, visitors.mobile, visitors.company_name, visitors.picture, visitors.sign_in, visitors.sign_out, visitors.employee_id, employees.fullname AS employee_name, employees.email AS employee_email
FROM visitors
LEFT JOIN employees ON visitors.employee_id=employees.id
WHERE sign_in >= NOW() - INTERVAL '24 HOURS'
AND sign_out IS NULL
ORDER by visitors.id
`

type ListVisitorsRow struct {
	ID            int64          `json:"id"`
	Fullname      string         `json:"fullname"`
	Email         string         `json:"email"`
	Mobile        string         `json:"mobile"`
	CompanyName   string         `json:"company_name"`
	Picture       string         `json:"picture"`
	SignIn        time.Time      `json:"sign_in"`
	SignOut       sql.NullTime   `json:"sign_out"`
	EmployeeID    sql.NullInt64  `json:"employee_id"`
	EmployeeName  sql.NullString `json:"employee_name"`
	EmployeeEmail sql.NullString `json:"employee_email"`
}

// SELECT * FROM visitors
// WHERE id = $1 LIMIT 1;
// -- name: ListVisitors :many
// SELECT * FROM visitors
// ORDER BY id
// LIMIT $1
// OFFSET $2;
func (q *Queries) ListVisitors(ctx context.Context) ([]ListVisitorsRow, error) {
	rows, err := q.db.QueryContext(ctx, listVisitors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListVisitorsRow{}
	for rows.Next() {
		var i ListVisitorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Fullname,
			&i.Email,
			&i.Mobile,
			&i.CompanyName,
			&i.Picture,
			&i.SignIn,
			&i.SignOut,
			&i.EmployeeID,
			&i.EmployeeName,
			&i.EmployeeEmail,
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

const updateVisitor = `-- name: UpdateVisitor :one
UPDATE visitors
SET
sign_out = now(),
updated_at = now()
WHERE id = $1
RETURNING id, fullname, email, mobile, company_name, picture, sign_in, sign_out, employee_id, created_at, updated_at
`

func (q *Queries) UpdateVisitor(ctx context.Context, id int64) (Visitor, error) {
	row := q.db.QueryRowContext(ctx, updateVisitor, id)
	var i Visitor
	err := row.Scan(
		&i.ID,
		&i.Fullname,
		&i.Email,
		&i.Mobile,
		&i.CompanyName,
		&i.Picture,
		&i.SignIn,
		&i.SignOut,
		&i.EmployeeID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
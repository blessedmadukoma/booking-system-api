// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: signedvisitors.sql

package db

import (
	"context"
	"database/sql"
)

const createSignedVisitor = `-- name: CreateSignedVisitor :one
INSERT INTO signedvisitors (
 visitor_id,
 visit_id,
 signed_out
) VALUES (
 $1, $2, $3
)
RETURNING id, visitor_id, signed_in, signed_out, visit_id, created_at, updated_at
`

type CreateSignedVisitorParams struct {
	VisitorID int64         `json:"visitor_id"`
	VisitID   sql.NullInt64 `json:"visit_id"`
	SignedOut sql.NullTime  `json:"signed_out"`
}

func (q *Queries) CreateSignedVisitor(ctx context.Context, arg CreateSignedVisitorParams) (Signedvisitor, error) {
	row := q.db.QueryRowContext(ctx, createSignedVisitor, arg.VisitorID, arg.VisitID, arg.SignedOut)
	var i Signedvisitor
	err := row.Scan(
		&i.ID,
		&i.VisitorID,
		&i.SignedIn,
		&i.SignedOut,
		&i.VisitID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listLatestSignedVisitor = `-- name: ListLatestSignedVisitor :one
SELECT id, visitor_id, signed_in, signed_out, visit_id, created_at, updated_at FROM signedvisitors 
WHERE visitor_id = $1 
AND visit_id = $2
ORDER BY created_at DESC
LIMIT 1
`

type ListLatestSignedVisitorParams struct {
	VisitorID int64         `json:"visitor_id"`
	VisitID   sql.NullInt64 `json:"visit_id"`
}

func (q *Queries) ListLatestSignedVisitor(ctx context.Context, arg ListLatestSignedVisitorParams) (Signedvisitor, error) {
	row := q.db.QueryRowContext(ctx, listLatestSignedVisitor, arg.VisitorID, arg.VisitID)
	var i Signedvisitor
	err := row.Scan(
		&i.ID,
		&i.VisitorID,
		&i.SignedIn,
		&i.SignedOut,
		&i.VisitID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listSignedVisitors = `-- name: ListSignedVisitors :many
SELECT id, visitor_id, signed_in, signed_out, visit_id, created_at, updated_at FROM signedvisitors ORDER BY id
`

func (q *Queries) ListSignedVisitors(ctx context.Context) ([]Signedvisitor, error) {
	rows, err := q.db.QueryContext(ctx, listSignedVisitors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Signedvisitor{}
	for rows.Next() {
		var i Signedvisitor
		if err := rows.Scan(
			&i.ID,
			&i.VisitorID,
			&i.SignedIn,
			&i.SignedOut,
			&i.VisitID,
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

const updateSignedVisitors = `-- name: UpdateSignedVisitors :one
UPDATE signedvisitors
SET
signed_out = $3,
updated_at = now()
WHERE visit_id = $2
AND id = $1
RETURNING id, visitor_id, signed_in, signed_out, visit_id, created_at, updated_at
`

type UpdateSignedVisitorsParams struct {
	ID        int64         `json:"id"`
	VisitID   sql.NullInt64 `json:"visit_id"`
	SignedOut sql.NullTime  `json:"signed_out"`
}

func (q *Queries) UpdateSignedVisitors(ctx context.Context, arg UpdateSignedVisitorsParams) (Signedvisitor, error) {
	row := q.db.QueryRowContext(ctx, updateSignedVisitors, arg.ID, arg.VisitID, arg.SignedOut)
	var i Signedvisitor
	err := row.Scan(
		&i.ID,
		&i.VisitorID,
		&i.SignedIn,
		&i.SignedOut,
		&i.VisitID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

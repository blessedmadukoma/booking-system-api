-- name: CreateSignedVisitor :one
INSERT INTO signedvisitors (
 visitor_id,
 visit_id,
 signed_out
) VALUES (
 $1, $2, $3
)
RETURNING *;

-- name: ListSignedVisitors :many
SELECT * FROM signedvisitors ORDER BY id;

-- name: ListLatestSignedVisitor :one
SELECT * FROM signedvisitors 
WHERE visitor_id = $1 
AND visit_id = $2
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateSignedVisitors :one
UPDATE signedvisitors
SET
signed_out = $3,
updated_at = now()
WHERE visit_id = $2
AND id = $1
RETURNING *;
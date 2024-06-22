// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: logged_out.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getLoggedOutToken = `-- name: GetLoggedOutToken :one
SELECT token
FROM logged_out
WHERE token = $1
`

func (q *Queries) GetLoggedOutToken(ctx context.Context, token string) (string, error) {
	row := q.db.QueryRow(ctx, getLoggedOutToken, token)
	err := row.Scan(&token)
	return token, err
}

const logout = `-- name: Logout :exec
INSERT INTO logged_out (token, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4)
RETURNING token, created_at, updated_at, user_id
`

type LogoutParams struct {
	Token     string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	UserID    pgtype.UUID
}

func (q *Queries) Logout(ctx context.Context, arg LogoutParams) error {
	_, err := q.db.Exec(ctx, logout,
		arg.Token,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
	)
	return err
}

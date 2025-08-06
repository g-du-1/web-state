package pagestate

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	conn, _ := pgx.Connect(ctx, connStr)

	return &Repository{
		conn: conn,
	}, nil
}

func (r Repository) CreatePagestate(ctx context.Context, pagestate Pagestate) (Pagestate, error) {
	err := r.conn.QueryRow(ctx,
		"INSERT INTO pagestates (url, scroll_pos) VALUES ($1, $2) RETURNING id",
		pagestate.Url, pagestate.ScrollPos).Scan(&pagestate.Id)

	return pagestate, err
}

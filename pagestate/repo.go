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

func (r Repository) SavePagestate(ctx context.Context, pagestate Pagestate) (Pagestate, error) {
	// TODO: Update if exists

	err := r.conn.QueryRow(ctx,
		"INSERT INTO pagestates (url, scroll_pos, visible_text) VALUES ($1, $2, $3) RETURNING id, created_at",
		pagestate.Url, pagestate.ScrollPos, pagestate.VisibleText).Scan(&pagestate.Id, &pagestate.CreatedAt)

	return pagestate, err
}

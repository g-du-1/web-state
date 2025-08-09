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

func (r Repository) GetPagestate(ctx context.Context, url string) (Pagestate, error) {
	var pagestate Pagestate

	err := r.conn.QueryRow(ctx, "SELECT id, url, scroll_pos, visible_text, created_at FROM pagestates WHERE url = $1", url).Scan(&pagestate.Id, &pagestate.Url, &pagestate.ScrollPos, &pagestate.VisibleText, &pagestate.CreatedAt)

	return pagestate, err
}

func (r Repository) GetAllPagestates(ctx context.Context) ([]Pagestate, error) {
	rows, _ := r.conn.Query(ctx, "SELECT id, url, scroll_pos, visible_text, created_at FROM pagestates ORDER BY created_at DESC")

	defer rows.Close()

	var pagestates []Pagestate

	for rows.Next() {
		var pagestate Pagestate

		rows.Scan(&pagestate.Id, &pagestate.Url, &pagestate.ScrollPos, &pagestate.VisibleText, &pagestate.CreatedAt)

		pagestates = append(pagestates, pagestate)
	}

	return pagestates, rows.Err()
}

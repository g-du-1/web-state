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

func (r Repository) GetAllPagestates(ctx context.Context) ([]Pagestate, error) {
	rows, err := r.conn.Query(ctx, "SELECT id, url, scroll_pos FROM pagestates ORDER BY id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pagestates []Pagestate

	for rows.Next() {
		var pagestate Pagestate

		err := rows.Scan(&pagestate.Id, &pagestate.Url, &pagestate.ScrollPos)

		if err != nil {
			return nil, err
		}

		pagestates = append(pagestates, pagestate)
	}

	return pagestates, rows.Err()
}

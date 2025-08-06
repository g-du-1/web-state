package pagestate

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, connStr)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

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

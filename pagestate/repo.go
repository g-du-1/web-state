package pagestate

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connStr)

	if err != nil {
		return nil, err
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (r Repository) SavePagestate(ctx context.Context, pagestate Pagestate) (Pagestate, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO pagestates (url, scroll_pos, visible_text)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (url) DO UPDATE SET scroll_pos = EXCLUDED.scroll_pos, visible_text = EXCLUDED.visible_text, updated_at = CURRENT_TIMESTAMP
		 RETURNING id, updated_at`,
		pagestate.Url, pagestate.ScrollPos, pagestate.VisibleText).Scan(&pagestate.Id, &pagestate.UpdatedAt)

	return pagestate, err
}

func (r Repository) GetPagestate(ctx context.Context, url string) (Pagestate, error) {
	var pagestate Pagestate

	err := r.pool.QueryRow(ctx, "SELECT id, url, scroll_pos, visible_text, updated_at FROM pagestates WHERE url = $1", url).Scan(&pagestate.Id, &pagestate.Url, &pagestate.ScrollPos, &pagestate.VisibleText, &pagestate.UpdatedAt)

	return pagestate, err
}

func (r Repository) GetAllPagestates(ctx context.Context) ([]Pagestate, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, url, scroll_pos, visible_text, updated_at FROM pagestates ORDER BY updated_at DESC")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pagestates []Pagestate

	for rows.Next() {
		var pagestate Pagestate

		err := rows.Scan(&pagestate.Id, &pagestate.Url, &pagestate.ScrollPos, &pagestate.VisibleText, &pagestate.UpdatedAt)

		if err != nil {
			return nil, err
		}

		pagestates = append(pagestates, pagestate)
	}

	return pagestates, rows.Err()
}

func (r Repository) DeleteAllPageStates(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM pagestates")

	return err
}

func (r Repository) Close() {
	r.pool.Close()
}

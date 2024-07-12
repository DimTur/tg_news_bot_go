package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/DimTur/tg_news_bot_go/internal/model"
	"github.com/jmoiron/sqlx"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with db: %w", err)
	}
	defer conn.Close()

	var sources []dbSource
	if err := conn.SelectContext(ctx, &sources, `SELECT * FROM sources`); err != nil {
		return nil, fmt.Errorf("failed to fetch sources: %w", err)
	}

	var mappedSources []model.Source
	for _, source := range sources {
		mappedSources = append(mappedSources, model.Source(source))
	}

	return mappedSources, nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with db: %w", err)
	}
	defer conn.Close()

	var source dbSource
	if err := conn.SelectContext(ctx, &source, `SELECT * FROM sources id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to fetch source with id %d: %w", id, err)
	}

	return (*model.Source)(&source), nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to connect with db: %w", err)
	}
	defer conn.Close()

	var id int64
	row := conn.QueryRowxContext(
		ctx,
		`INSERT INTO sources (name, feed_url, created_at) VALUES ($1, $2, $3) RETURNING id`,
		source.Name,
		source.FeedURL,
		source.CreatedAt,
	)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("failed to add new source: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to scan source id: %w", err)
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect with db: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `DELETE FROM sources WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete source with id %d: %w", id, err)
	}

	return nil
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

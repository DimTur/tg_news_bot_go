package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/DimTur/tg_news_bot_go/internal/model"
	"github.com/jmoiron/sqlx"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (a *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect with db: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at)
		VALUES ($1, $2, $3, $4, $5
		ON CONFLICT DO NOTHING`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PostedAt,
	); err != nil {
		return fmt.Errorf("could not insert article: %w", err)
	}

	return nil
}

func (a *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with db: %w", err)
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles
		 WHERE posted_at IS NULL
			AND publiched_at >= $1::timestamp
		 ORDER BY published_at DESC 
		 LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit,
	); err != nil {
		return nil, fmt.Errorf("could not select articles: %w", err)
	}

	var mappedArticles []model.Article
	for _, article := range articles {
		mappedArticles = append(mappedArticles, model.Article(article))
	}

	return mappedArticles, nil
}

func (a *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return fmt.Errorf("could not get connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2`,
		time.Now().UTC().Format(time.RFC3339),
		id,
	); err != nil {
		return fmt.Errorf("could not mark article as posted: %w", err)
	}

	return nil
}

type dbArticle struct {
	ID          int64     `db:"id"`
	SourceID    int64     `db:"source_id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Summary     string    `db:"summary"`
	PublishedAt time.Time `db:"published_at"`
	PostedAt    time.Time `db:"posted_at"`
	CreatedAt   time.Time `db:"created_at"`
}

package repository

import (
	"context"
	"database/sql"

	"calanggo-server/internal/core/domain"
	"calanggo-server/internal/core/ports"

	// Importa o driver mas não o usa diretamente (apenas registra no database/sql)
	_ "github.com/glebarez/go-sqlite"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dsn string) (ports.LinkRepository, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS links (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		shortened TEXT NOT NULL,
		created_at DATETIME,
		visits INTEGER DEFAULT 0
	);
	`
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &sqliteRepository{db: db}, nil
}

func (r *sqliteRepository) Save(ctx context.Context, link *domain.Link) error {
	query := `
		INSERT INTO links (id, original_url, shortened, created_at, visits)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		link.ID,
		link.Original,
		link.Shortened,
		link.CreatedAt,
		link.Visits,
	)
	return err
}

func (r *sqliteRepository) GetByShortCode(ctx context.Context, code string) (*domain.Link, error) {
	query := `
		SELECT id, original_url, shortened, created_at, visits
		FROM links WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, code)

	var link domain.Link
	err := row.Scan(
		&link.ID,
		&link.Original,
		&link.Shortened,
		&link.CreatedAt,
		&link.Visits,
	)

	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *sqliteRepository) IncrementVisits(ctx context.Context, code string) error {
	query := `UPDATE links SET visits = visits + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, code)
	return err
}

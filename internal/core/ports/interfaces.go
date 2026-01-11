package ports

import (
	"calanggo-server/internal/core/domain"
	"context"
)

type LinkRepository interface {
	Save(ctx context.Context, link *domain.Link) error
	GetByShortCode(ctx context.Context, code string) (*domain.Link, error)
	IncrementVisits(ctx context.Context, code string) error
}

type LinkService interface {
	CreateShortLink(ctx context.Context, originalURL string) (*domain.Link, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}

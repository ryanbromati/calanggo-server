package services

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"calanggo-server/internal/core/domain"
	"calanggo-server/internal/core/ports"
	"calanggo-server/pkg/base62"
)

type linkService struct {
	repo ports.LinkRepository
}

func NewLinkService(repo ports.LinkRepository) ports.LinkService {
	return &linkService{
		repo: repo,
	}
}

func (s *linkService) CreateShortLink(ctx context.Context, originalURL string) (*domain.Link, error) {
	if originalURL == "" {
		return nil, errors.New("URL original não pode ser vazia")
	}

	// TODO: Em produção, isso viria de um gerador de ID distribuído (Snowflake) ou do banco (RETURNING id)
	// gerando um ID aleatório para uso local
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := rnd.Uint64()

	shortCode := base62.Encode(id)

	link := domain.NewLink(originalURL, shortCode)

	if err := s.repo.Save(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *linkService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	link, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	_ = s.repo.IncrementVisits(ctx, shortCode)

	return link.Original, nil
}

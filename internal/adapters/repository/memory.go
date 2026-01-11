package repository

import (
	"context"
	"errors"
	"sync"

	"calanggo-server/internal/core/domain"
	"calanggo-server/internal/core/ports"
)

type memoryRepository struct {
	links    map[string]*domain.Link
	mu       sync.RWMutex
	visitsCh chan string
}

func NewMemoryRepository() ports.LinkRepository {
	repo := &memoryRepository{
		links:    make(map[string]*domain.Link),
		visitsCh: make(chan string, 100),
	}

	go repo.backgroundVisitWorker()

	return repo
}

func (r *memoryRepository) Save(ctx context.Context, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.links[link.ID]; exists {
		return errors.New("link já existe")
	}

	r.links[link.ID] = link
	return nil
}

func (r *memoryRepository) GetByShortCode(ctx context.Context, code string) (*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, exists := r.links[code]
	if !exists {
		return nil, errors.New("link não encontrado")
	}
	return link, nil
}

func (r *memoryRepository) IncrementVisits(ctx context.Context, code string) error {
	select {
	case r.visitsCh <- code:
		return nil
	default:
		return errors.New("fila de visitas cheia")
	}
}

func (r *memoryRepository) backgroundVisitWorker() {
	for code := range r.visitsCh {
		r.mu.Lock()
		if link, exists := r.links[code]; exists {
			link.Visits++
		}
		r.mu.Unlock()
	}
}

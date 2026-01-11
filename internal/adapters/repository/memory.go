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
	visitsCh chan string // Canal assíncrono
}

// NewMemoryRepository cria um repositório e inicia o worker de background
func NewMemoryRepository() ports.LinkRepository {
	repo := &memoryRepository{
		links:    make(map[string]*domain.Link),
		visitsCh: make(chan string, 100), // Buffer para não bloquear o sender imediatamente
	}

	// Inicia o worker que processa as atualizações
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

// IncrementVisits agora é Non-Blocking (apenas envia pro canal)
func (r *memoryRepository) IncrementVisits(ctx context.Context, code string) error {
	// Select com default impede que o servidor trave se o canal encher (backpressure simples)
	select {
	case r.visitsCh <- code:
		return nil
	default:
		// Em produção, logariamos que o buffer encheu e perdemos uma contagem
		return errors.New("fila de visitas cheia")
	}
}

// backgroundVisitWorker consome o canal e aplica as mudanças
func (r *memoryRepository) backgroundVisitWorker() {
	// range no canal lê continuamente até o canal ser fechado
	for code := range r.visitsCh {
		r.mu.Lock()
		if link, exists := r.links[code]; exists {
			link.Visits++
		}
		r.mu.Unlock()
	}
}

package oauth

import (
	"sync"
	"user-review-ingest/internal/domain/errors"
)

type ProviderRegistry struct {
	providers map[string]Provider
	mutex     sync.RWMutex
}

func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

func (r *ProviderRegistry) Register(name string, provider Provider) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.providers[name] = provider
}

func (r *ProviderRegistry) GetProvider(name string) (Provider, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, errors.ErrInvalidProvider
	}

	return provider, nil
}

func (r *ProviderRegistry) ListProviders() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	providers := make([]string, 0, len(r.providers))
	for name := range r.providers {
		providers = append(providers, name)
	}
	return providers
}

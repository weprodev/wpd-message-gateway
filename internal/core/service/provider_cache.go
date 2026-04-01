package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/registry"
)

// providerCache caches instantiated send providers keyed by integration ID and
// the update timestamp. The timestamp acts as a natural version: when the
// integration config changes, its updated_at advances and the old entry is
// ignored on the next lookup, so the factory is called exactly once per
// effective config version.
//
// The cache is safe for concurrent use. Factory construction happens outside
// the write lock to avoid holding it during potentially slow HTTP-client setup.
type providerCache[F any] struct {
	mu      sync.RWMutex
	entries map[cacheKey]F
}

type cacheKey struct {
	id        string
	updatedAt int64 // integration.UpdatedAt.UnixNano()
}

func newProviderCache[F any]() *providerCache[F] {
	return &providerCache[F]{entries: make(map[cacheKey]F)}
}

// get returns a cached provider, true if found.
func (c *providerCache[F]) get(intg *domain.Integration) (F, bool) {
	k := cacheKey{intg.ID, intg.UpdatedAt.UnixNano()}
	c.mu.RLock()
	f, ok := c.entries[k]
	c.mu.RUnlock()
	return f, ok
}

// set stores a provider in the cache.
func (c *providerCache[F]) set(intg *domain.Integration, f F) {
	k := cacheKey{intg.ID, intg.UpdatedAt.UnixNano()}
	c.mu.Lock()
	c.entries[k] = f
	c.mu.Unlock()
}

// emailSenderCache resolves (and caches) port.EmailSenders from integrations.
type emailSenderCache = providerCache[port.EmailSender]

// smsSenderCache resolves (and caches) port.SMSSenders from integrations.
type smsSenderCache = providerCache[port.SMSSender]

// pushSenderCache resolves (and caches) port.PushSenders from integrations.
type pushSenderCache = providerCache[port.PushSender]

// chatSenderCache resolves (and caches) port.ChatSenders from integrations.
type chatSenderCache = providerCache[port.ChatSender]

// resolveEmailSender returns a cached or newly constructed EmailSender for the integration.
func resolveEmailSender(cache *emailSenderCache, intg *domain.Integration) (port.EmailSender, error) {
	if s, ok := cache.get(intg); ok {
		return s, nil
	}
	factory, err := registry.GetEmailFactory(intg.ProviderName)
	if err != nil {
		return nil, fmt.Errorf("provider factory: %w", err)
	}
	var cfg registry.EmailConfig
	if err := json.Unmarshal(intg.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parse integration config: %w", err)
	}
	sender, err := factory(cfg)
	if err != nil {
		return nil, fmt.Errorf("init email provider %q: %w", intg.ProviderName, err)
	}
	cache.set(intg, sender)
	return sender, nil
}

// resolveSMSSender returns a cached or newly constructed SMSSender.
func resolveSMSSender(cache *smsSenderCache, intg *domain.Integration) (port.SMSSender, error) {
	if s, ok := cache.get(intg); ok {
		return s, nil
	}
	factory, err := registry.GetSMSFactory(intg.ProviderName)
	if err != nil {
		return nil, fmt.Errorf("provider factory: %w", err)
	}
	var cfg registry.SMSConfig
	if err := json.Unmarshal(intg.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parse integration config: %w", err)
	}
	sender, err := factory(cfg)
	if err != nil {
		return nil, fmt.Errorf("init SMS provider %q: %w", intg.ProviderName, err)
	}
	cache.set(intg, sender)
	return sender, nil
}

// resolvePushSender returns a cached or newly constructed PushSender.
func resolvePushSender(cache *pushSenderCache, intg *domain.Integration) (port.PushSender, error) {
	if s, ok := cache.get(intg); ok {
		return s, nil
	}
	factory, err := registry.GetPushFactory(intg.ProviderName)
	if err != nil {
		return nil, fmt.Errorf("provider factory: %w", err)
	}
	var cfg registry.PushConfig
	if err := json.Unmarshal(intg.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parse integration config: %w", err)
	}
	sender, err := factory(cfg)
	if err != nil {
		return nil, fmt.Errorf("init push provider %q: %w", intg.ProviderName, err)
	}
	cache.set(intg, sender)
	return sender, nil
}

// resolveChatSender returns a cached or newly constructed ChatSender.
func resolveChatSender(cache *chatSenderCache, intg *domain.Integration) (port.ChatSender, error) {
	if s, ok := cache.get(intg); ok {
		return s, nil
	}
	factory, err := registry.GetChatFactory(intg.ProviderName)
	if err != nil {
		return nil, fmt.Errorf("provider factory: %w", err)
	}
	var cfg registry.ChatConfig
	if err := json.Unmarshal(intg.Config, &cfg); err != nil {
		return nil, fmt.Errorf("parse integration config: %w", err)
	}
	sender, err := factory(cfg)
	if err != nil {
		return nil, fmt.Errorf("init chat provider %q: %w", intg.ProviderName, err)
	}
	cache.set(intg, sender)
	return sender, nil
}

// ensure time is used (for cacheKey.updatedAt type alignment).
var _ = time.Now

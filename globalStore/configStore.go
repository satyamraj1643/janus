package configStore

import (
	"encoding/json"
	"sync"
)

type CachedConfig struct {
	Config json.RawMessage
	ID     string
}

type store struct {
	mu     sync.RWMutex
	config map[string]CachedConfig
}

var s = &store{
	config: make(map[string]CachedConfig),
}

func Get(userID string) (CachedConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg, ok := s.config[userID]
	return cfg, ok
}

func Set(userID string, cfg json.RawMessage, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config[userID] = CachedConfig{Config: cfg, ID: id}
}

func Delete(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.config, userID)
}

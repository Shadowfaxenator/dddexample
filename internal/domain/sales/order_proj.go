package sales

import (
	"context"
	"sync"

	"github.com/alekseev-bro/ddd/pkg/domain"
)

type DB struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewRamDB() *DB {
	return &DB{
		data: make(map[string]any),
	}
}

func (db *DB) Get(key string) (any, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	val, ok := db.data[key]
	return val, ok
}

func (db *DB) Set(key string, val any) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = val
}

type OrderProjection struct {
	db *DB
}

func (c *OrderProjection) Handle(ctx context.Context, eventID domain.EventID[Order], e domain.Event[Order]) error {
	switch ev := e.(type) {
	case *OrderCreated:
		c.db.Set(string(eventID), ev.Order)
	}
	return nil
}

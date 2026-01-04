package sales

import (
	"context"
	"sync"

	"github.com/alekseev-bro/ddd/pkg/eventstore"
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

type OrderSaga struct {
	cust eventstore.EventStore[Customer]
}

func (c *OrderSaga) Handle(ctx context.Context, o *Order, eventID eventstore.EventID[Order]) error {

	return nil
}

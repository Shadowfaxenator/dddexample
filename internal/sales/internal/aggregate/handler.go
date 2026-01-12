package aggregate

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"
)

type EventHandler[T any] interface {
	Handle(ctx context.Context, eventID string, event T) error
}

type CommandHandler[T, U any] interface {
	Handle(ctx context.Context, id events.ID[T], cmd U, idempotencyKey string) error
}

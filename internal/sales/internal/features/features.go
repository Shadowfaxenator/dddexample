package features

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"
)

type CommandHandler[T, U any] interface {
	Handle(ctx context.Context, id events.ID[T], cmd U, idempotencyKey string) error
}

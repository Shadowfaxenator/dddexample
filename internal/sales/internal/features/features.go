package features

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
)

type CommandHandler[T, U any] interface {
	Handle(ctx context.Context, id essrv.ID[T], cmd U, idempotencyKey string) error
}

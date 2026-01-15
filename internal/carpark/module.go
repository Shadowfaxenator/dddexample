package carpark

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type Module struct {
	aggregate.Aggregate
	// ...
}

type handler struct {
}

func (h *handler) HandleEvent(ctx context.Context, eventID string) error {

	// Handle the event
	return nil
}

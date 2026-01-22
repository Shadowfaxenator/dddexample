package carpark

import (
	"context"
)

type Module struct {

	// ...
}

type handler struct {
}

func (h *handler) HandleEvent(ctx context.Context, eventID string) error {

	// Handle the event
	return nil
}

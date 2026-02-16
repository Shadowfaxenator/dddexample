package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/stream"
	"github.com/alekseev-bro/dddexample/internal/carpark/internal/aggregate/car"
)

type RegisterCar struct {
	Car *car.Car
}

type registerCarHandler struct {
	Cars aggregate.Mutator[car.Car, *car.Car]
}

func NewRegisterCarHandler(cars aggregate.Mutator[car.Car, *car.Car]) *registerCarHandler {
	return &registerCarHandler{Cars: cars}
}

func (h *registerCarHandler) HandleCommand(ctx context.Context, cmd RegisterCar) ([]stream.EventMetadata, error) {
	return h.Cars.Mutate(ctx, cmd.Car.ID, func(state *car.Car) (aggregate.Events[car.Car], error) {
		return state.Register(cmd.Car)
	})
}

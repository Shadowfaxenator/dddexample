package carpark

import (
	"context"
	"ddd/pkg/aggregate"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Service) CreateCar(ctx context.Context, model string, brand string) error {
	id := aggregate.NewID()
	return s.car.CommandFunc(ctx, id, func(c *Car) (*aggregate.Event[Car], error) {

		event := aggregate.NewEvent(CarCreated{
			Car: Car{VIN: uuid.New().String(),
				CarModel: CarModel{Brand: brand, Model: model}, RentState: Available, MaintananceState: NotNeeded},
		})
		return event, nil

	})

}
func (s *Service) RentCar(ctx context.Context, orderID aggregate.ID, carID aggregate.ID) error {
	return s.car.CommandFunc(ctx, carID, func(c *Car) (*aggregate.Event[Car], error) {
		if c.RentState == Available {

			return aggregate.NewEvent(CarRented{OrderID: orderID}), nil
		}
		slog.Warn("Car rent rejected")
		return aggregate.NewEvent(CarRentRejected{OrderID: orderID}), nil
	})
}

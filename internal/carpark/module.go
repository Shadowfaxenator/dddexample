package carpark

import (
	"context"
	"log/slog"
	"os"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/codec"
	"github.com/alekseev-bro/ddd/pkg/natsaggregate"
	"github.com/alekseev-bro/dddexample/internal/carpark/internal/aggregate/car"
	carcmd "github.com/alekseev-bro/dddexample/internal/carpark/internal/aggregate/car/command"
	"github.com/alekseev-bro/dddexample/internal/carpark/internal/integration"
	"github.com/nats-io/nats.go/jetstream"
)

type Module struct {
	RegisterCarHandler aggregate.CommandHandler[car.Car, carcmd.RegisterCar]
}

func NewModule(ctx context.Context, js jetstream.JetStream, publisher integration.Publisher) *Module {
	cars, err := natsaggregate.New(ctx, js,
		natsaggregate.WithInMemory[car.Car](),
		natsaggregate.WithEvent[car.Arrived, car.Car](),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	d, err := cars.Subscribe(ctx, integration.NewCarHandler(publisher, codec.JSON))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()
		d.Drain()
	}()

	return &Module{
		RegisterCarHandler: carcmd.NewRegisterCarHandler(cars),
	}
}

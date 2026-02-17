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
	RegisterCarHandler aggregate.CommandHandler[carcmd.RegisterCar, car.Car]
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
	if err = cars.Subscribe(ctx, integration.NewCarHandler(publisher, codec.JSON)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()
		if err := cars.Drain(); err != nil {
			slog.Error(err.Error())
		}
	}()

	return &Module{
		RegisterCarHandler: carcmd.NewRegisterCarHandler(cars),
	}
}

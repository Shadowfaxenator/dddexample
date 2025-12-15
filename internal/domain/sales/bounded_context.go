package sales

import (
	"context"
	"log/slog"
	"time"

	"github.com/alekseev-bro/ddd/pkg/saga"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"

	"github.com/alekseev-bro/ddd/pkg/domain"

	"github.com/nats-io/nats.go/jetstream"
)

type boundedContext struct {
	Customer domain.Aggregate[Customer]
	Order    domain.Aggregate[Order]
}

// type MySerder struct {
// }

// func (m *MySerder) Serialize(in any) ([]byte, error) {

// 	var buf bytes.Buffer
// 	if err := gob.NewEncoder(&buf).Encode(in); err != nil {

// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func (m *MySerder) Deserialize(data []byte, out any) error {

// 	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(out); err != nil {
// 		return err
// 	}
// 	return nil
// }

func New(ctx context.Context, js jetstream.JetStream) *boundedContext {

	customer := domain.NewAggregate(ctx,
		esnats.NewEventStream(ctx, js, esnats.WithInMemory[Customer]()),
		snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Customer]()),
		domain.WithSnapshotThreshold[Customer](10, time.Second),
	)

	domain.RegisterEvent[*OrderRejected](customer)
	domain.RegisterEvent[*CustomerCreated](customer)
	domain.RegisterEvent[*OrderAccepted](customer)
	oes := esnats.NewEventStream(ctx, js, esnats.WithInMemory[Order]())
	snap := snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Order]())

	order := domain.NewAggregate(ctx, oes, snap, domain.WithSnapshotThreshold[Order](10, time.Second))

	domain.RegisterEvent[*OrderCreated](order)
	domain.RegisterEvent[*OrderClosed](order)
	domain.RegisterEvent[*OrderVerified](order)

	var subs []domain.Drainer
	sub, err := order.Project(ctx, &OrderProjection{
		db: NewRamDB(),
	})
	if err != nil {
		panic(err)
	}

	subs = append(subs, sub)

	saga.Step(ctx, order, customer, func(e *OrderCreated) *ValidateOrder {

		return &ValidateOrder{CustomerID: e.Order.CustomerID, OrderID: e.Order.ID}
	})

	go func() {
		<-ctx.Done()
		for _, sub := range subs {
			sub.Drain()

		}
		slog.Info("all subscriptions closed")

	}()

	return &boundedContext{
		Customer: customer,
		Order:    order,
	}
}

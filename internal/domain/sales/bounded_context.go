package sales

import (
	"context"
	"log/slog"
	"time"

	"github.com/alekseev-bro/ddd/pkg/domain"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"

	"github.com/alekseev-bro/ddd/pkg/aggregate"

	"github.com/nats-io/nats.go/jetstream"
)

type boundedContext struct {
	Customer aggregate.Aggregate[Customer]
	Order    aggregate.Aggregate[Order]
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
//
//

func New(ctx context.Context, js jetstream.JetStream) *boundedContext {

	customer := domain.NewAggregate(ctx,
		esnats.NewEventStream(ctx, js, esnats.WithInMemory[Customer]()),
		snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Customer]()),
		aggregate.WithSnapshotThreshold[Customer](10, time.Second),
		domain.WithEvent[*OrderRejected](),
		domain.WithEvent[*OrderAccepted](),
	)

	order := natsstore.NewAggregate(ctx, js,
		natsstore.WithSnapshotThreshold[Order](10, time.Second),
		natsstore.WithInMemory[Order](),
		natsstore.WithEvent[*OrderClosed](),
	)

	bc := &boundedContext{
		Customer: customer,
		Order:    order,
	}

	return bc
}

func (b *boundedContext) StartOrderCreationSaga(ctx context.Context) {
	domain.SagaStep(ctx, b.Order, b.Customer, func(e *domain.Created[Order]) *ValidateOrder {
		return &ValidateOrder{CustomerID: e.Body.CustomerID, OrderID: e.ID}
	})
	domain.SagaStep(ctx, b.Order, b.Customer, func(e *domain.Created[Order]) *domain.Create[Customer] {

		return &domain.Create[Customer]{ID: b.Customer.NewID(), Body: &Customer{Name: "ddd"}}
	})
}

func (b *boundedContext) StartProjections(ctx context.Context) {
	var subs []aggregate.Drainer
	sub, err := b.Order.Project(ctx, &OrderProjection{
		db: NewRamDB(),
	})
	if err != nil {
		panic(err)
	}

	subs = append(subs, sub)

	go func() {
		<-ctx.Done()
		for _, sub := range subs {
			sub.Drain()

		}
		slog.Info("all subscriptions closed")

	}()
}

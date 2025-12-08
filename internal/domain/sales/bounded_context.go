package sales

import (
	"context"
	"ddd/pkg/domain"
	"ddd/pkg/store/natsstore/esnats"
	"ddd/pkg/store/natsstore/snapnats"
	"log/slog"

	"github.com/nats-io/nats.go"
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

func New(ctx context.Context) *boundedContext {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	customer := domain.NewAggregate[Customer](ctx,
		esnats.NewEventStream(ctx, js, esnats.WithInMemory[Customer]()),
		snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Customer]()),
	)

	domain.RegisterEvent[*OrderRejected](customer)
	domain.RegisterEvent[*CustomerCreated](customer)
	domain.RegisterEvent[*OrderAccepted](customer)
	oes := esnats.NewEventStream(ctx, js, esnats.WithInMemory[Order]())
	snap := snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Order]())

	order := domain.NewAggregate[Order](ctx, oes, snap)

	domain.RegisterEvent[*OrderCreated](order)
	domain.RegisterEvent[*OrderClosed](order)
	domain.RegisterEvent[*OrderVerified](order)

	order.Project(ctx, &OrderService{
		Customer: customer,
	})
	var subs []domain.Drainer

	sub, _ := customer.Project(ctx, &CustomerService{
		Order: order,
	}, domain.FilterByEvent[*OrderAccepted]())

	subs = append(subs, sub...)

	ss := domain.Saga(ctx, order, customer, func(e *OrderCreated) *ValidateOrder {
		return &ValidateOrder{CustomerID: e.Order.CustomerID, OrderID: e.Order.ID}
	})

	subs = append(subs, ss)

	go func() {
		<-ctx.Done()
		for _, sub := range subs {
			sub.Drain()

		}
		slog.Info("all subscriptions drained")
	}()

	return &boundedContext{
		Customer: customer,
		Order:    order,
	}
}

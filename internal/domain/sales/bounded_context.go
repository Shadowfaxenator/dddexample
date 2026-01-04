package sales

import (
	"context"
	"time"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"

	"github.com/alekseev-bro/ddd/pkg/eventstore"

	"github.com/nats-io/nats.go/jetstream"
)

type boundedContext struct {
	Customer eventstore.EventStore[Customer]
	Order    eventstore.EventStore[Order]
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

	customer := eventstore.New(ctx,
		esnats.NewEventStream(ctx, js, esnats.WithInMemory[Customer]()),
		snapnats.NewSnapshotStore(ctx, js, snapnats.WithInMemory[Customer]()),
		eventstore.WithSnapshotThreshold[Customer](10, time.Second),
		eventstore.WithEvent[OrderRejected](),
		eventstore.WithEvent[OrderAccepted](),
		eventstore.WithEvent[CustomerCreated](),
	)

	order := natsstore.NewAggregate(ctx, js,
		natsstore.WithSnapshotThreshold[Order](10, time.Second),
		natsstore.WithInMemory[Order](),
		natsstore.WithEvent[OrderClosed](),
		natsstore.WithEvent[OrderCreated](),
	)

	bc := &boundedContext{
		Customer: customer,
		Order:    order,
	}

	return bc
}

func (b *boundedContext) StartOrderCreationSaga(ctx context.Context) {
	b.Customer.ProjectEvent(ctx, CustomerService{Order: b.Order})
	// aggregate.SagaStep(ctx, b.Order, b.Customer, func(e *event.Created[Order]) *ValidateOrder {
	// 	return &ValidateOrder{CustomerID: e.Body.CustomerID, OrderID: e.ID}
	// })

}

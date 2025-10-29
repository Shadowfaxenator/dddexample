package sales

import (
	"bytes"
	"context"
	"ddd/pkg/aggregate"
	"ddd/pkg/store/esnats"
	"encoding/gob"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Commander[T any] interface {
	Command(ctx context.Context, id aggregate.ID[T], command aggregate.Command[T]) error
}
type Subscriber[T any] interface {
	Subscribe(ctx context.Context, name string, handler func(aggregate.Event[T]) error, ordered bool)
}

type Aggregate[T any] interface {
	Commander[T]
	Subscriber[T]
}

type boundedContext struct {
	Customer        Aggregate[Customer]
	Order           Aggregate[Order]
	orderService    *OrderService
	customerService *CustomerService
}

type MySerder struct {
}

func (m *MySerder) Serialize(in any) ([]byte, error) {

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {

		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *MySerder) Deserialize(data []byte, out any) error {

	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(out); err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context) *boundedContext {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}
	custStream := esnats.NewEventStream[Customer](ctx, js)
	customer := aggregate.New(ctx,
		custStream,
		esnats.NewSnapshotStore[Customer](ctx, js),
	)
	aggregate.RegisterEvent[CustomerCreated](customer)
	//gob.Register(CustomerCreated{})
	// gob.Register(OrderAccepted{})
	// gob.Register(OrderCreated{})
	// gob.Register(OrderClosed{})
	// gob.Register(OrderVerified{})
	aggregate.RegisterEvent[OrderAccepted](customer)
	//aggregate.RegisterCommand[CreateCustomer](customer)

	order := aggregate.New[Order](ctx,
		esnats.NewEventStream[Order](ctx, js),
		esnats.NewSnapshotStore[Order](ctx, js),
	)
	aggregate.RegisterEvent[OrderCreated](order)
	aggregate.RegisterEvent[OrderClosed](order)
	aggregate.RegisterEvent[OrderVerified](order)
	//aggregate.RegisterCommand[CreateOrder](order)
	c := &boundedContext{
		Customer:        customer,
		Order:           order,
		orderService:    NewOrderService(ctx, customer, order),
		customerService: NewCustomerService(ctx, customer, order),
	}
	return c
}

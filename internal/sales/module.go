package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/orders"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"

	"github.com/nats-io/nats.go/jetstream"
)

type EventHandler[T any] interface {
	Handle(ctx context.Context, eventID string, event T) error
}

type Module struct {
	OrderPostedHandler EventHandler[orders.OrderPosted]
	RegisterCustomer   features.CommandHandler[customers.Customer, customer.Register]
	PostOrder          features.CommandHandler[orders.Order, order.Post]
	OrderStream        essrv.Projector[orders.Order]
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {

	cust := essrv.New(ctx,
		esnats.NewEventStream[customers.Customer](ctx, js, esnats.EventStreamConfig{
			StoreType: esnats.Memory,
		}),
		snapnats.NewSnapshotStore[customers.Customer](ctx, js, snapnats.SnapshotStoreConfig{
			StoreType: snapnats.Memory,
		}),
		essrv.AggregateConfig{
			SnapthotMsgThreshold: 5,
		},
		essrv.WithEvent[customers.OrderRejected](),
		essrv.WithEvent[customers.OrderAccepted](),
		essrv.WithEvent[customers.CustomerRegistered](),
	)

	ord := natsstore.NewAggregate(ctx, js,
		natsstore.NatsAggregateConfig{
			AggregateConfig: essrv.AggregateConfig{
				SnapthotMsgThreshold: 5,
			},
			StoreType: natsstore.Memory,
		},
		essrv.WithEvent[orders.OrderClosed](),
		essrv.WithEvent[orders.OrderPosted](),
		essrv.WithEvent[orders.OrderVerified](),
	)
	mod := &Module{
		PostOrder:          order.NewPostOrderHandler(ord),
		RegisterCustomer:   customer.NewRegisterHandler(cust),
		OrderPostedHandler: customer.NewOrderPostedHandler(customer.NewVerifyOrderHandler(cust)),
		OrderStream:        ord,
	}

	return mod
}

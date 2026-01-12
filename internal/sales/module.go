package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	customercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/command"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	ordercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/command"

	"github.com/nats-io/nats.go/jetstream"
)

type Module struct {
	OrderPostedHandler events.EventHandler[order.Posted]
	RegisterCustomer   aggregate.CommandHandler[customer.Customer, customercmd.Register]
	PostOrder          aggregate.CommandHandler[order.Order, ordercmd.Post]
	OrderStream        events.Subscriber[order.Order]
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {

	cust := events.NewStore(ctx,
		esnats.NewEventStream[customer.Customer](ctx, js, esnats.EventStreamConfig{
			StoreType: esnats.Memory,
		}),
		snapnats.NewSnapshotStore[customer.Customer](ctx, js, snapnats.SnapshotStoreConfig{
			StoreType: snapnats.Memory,
		}),
		events.AggregateConfig{
			SnapthotMsgThreshold: 5,
		},

		events.WithEvent[customer.OrderRejected](),
		events.WithEvent[customer.OrderAccepted](),
		events.WithEvent[customer.Registered](),
	)

	ord := natsstore.NewStore(ctx, js,
		natsstore.NatsAggregateConfig{
			AggregateConfig: events.AggregateConfig{
				SnapthotMsgThreshold: 5,
			},
			StoreType: natsstore.Memory,
		},
		events.WithEvent[order.Closed](),
		events.WithEvent[order.Posted](),
		events.WithEvent[order.Verified](),
	)

	mod := &Module{
		PostOrder:          ordercmd.NewPostOrderHandler(ord),
		RegisterCustomer:   customercmd.NewRegisterHandler(cust),
		OrderPostedHandler: customercmd.NewOrderPostedHandler(customercmd.NewVerifyOrderHandler(cust)),
		OrderStream:        ord,
	}

	return mod
}

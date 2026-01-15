package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	customercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/command"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	ordercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/command"

	"github.com/nats-io/nats.go/jetstream"
)

type Module struct {
	RegisterCustomer aggregate.CommandHandler[customer.Customer, customercmd.Register]
	PostOrder        aggregate.CommandHandler[order.Order, ordercmd.Post]
	OrderStream      aggregate.Subscriber[order.Order]
	CustomerStream   aggregate.Subscriber[customer.Customer]
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {

	cust := aggregate.NewStore(ctx,
		esnats.NewEventStream[customer.Customer](ctx, js, esnats.EventStreamConfig{
			StoreType: esnats.Memory,
		}),
		snapnats.NewSnapshotStore[customer.Customer](ctx, js, snapnats.SnapshotStoreConfig{
			StoreType: snapnats.Memory,
		}),
		aggregate.AggregateConfig{
			SnapthotMsgThreshold: 5,
		},

		aggregate.WithEvent[customer.OrderRejected]("OrderRejected"),
		aggregate.WithEvent[customer.OrderAccepted]("OrderAccepted"),
		aggregate.WithEvent[customer.Registered]("CustomerRegistered"),
	)

	ord := natsstore.NewStore(ctx, js,
		natsstore.NatsAggregateConfig{
			AggregateConfig: aggregate.AggregateConfig{
				SnapthotMsgThreshold: 5,
			},
			StoreType: natsstore.Memory,
		},
		aggregate.WithEvent[order.Closed]("OrderClosed"),
		aggregate.WithEvent[order.Posted]("OrderPosted"),
		aggregate.WithEvent[order.Verified]("OrderVerified"),
	)
	aggregate.Project(ctx, ord, customercmd.NewOrderPostedHandler(
		customercmd.NewVerifyOrderHandler(cust),
	))
	aggregate.Project(ctx, cust, ordercmd.NewOrderRejectedHandler(
		ordercmd.NewCloseOrderHandler(ord),
	))

	mod := &Module{
		PostOrder:        ordercmd.NewPostOrderHandler(ord),
		RegisterCustomer: customercmd.NewRegisterHandler(cust),
		OrderStream:      ord,
		CustomerStream:   cust,
	}

	return mod
}

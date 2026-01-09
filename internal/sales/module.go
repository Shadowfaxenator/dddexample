package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/esnats"
	"github.com/alekseev-bro/ddd/pkg/store/natsstore/snapnats"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer"
	customerUsecase "github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer/usecase"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
	orderUsecase "github.com/alekseev-bro/dddexample/internal/sales/internal/features/order/usecase"

	"github.com/nats-io/nats.go/jetstream"
)

type EventHandler[T any] interface {
	Handle(ctx context.Context, eventID string, event T) error
}

type Module struct {
	OrderPostedHandler EventHandler[order.Posted]
	RegisterCustomer   features.CommandHandler[customer.Customer, customerUsecase.Register]
	PostOrder          features.CommandHandler[order.Order, orderUsecase.Post]
	OrderStream        events.Projector[order.Order]
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
		PostOrder:          orderUsecase.NewPostOrderHandler(ord),
		RegisterCustomer:   customerUsecase.NewRegisterHandler(cust),
		OrderPostedHandler: customerUsecase.NewOrderPostedHandler(customerUsecase.NewVerifyOrderHandler(cust)),
		OrderStream:        ord,
	}

	return mod
}

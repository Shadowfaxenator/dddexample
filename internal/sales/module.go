package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
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
	RegisterCustomer   features.CommandHandler[customer.AggregateRoot, customerUsecase.Register]
	PostOrder          features.CommandHandler[order.AggregateRoot, orderUsecase.Post]
	OrderStream        essrv.Projector[order.AggregateRoot]
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {

	cust := essrv.New(ctx,
		esnats.NewEventStream[customer.AggregateRoot](ctx, js, esnats.EventStreamConfig{
			StoreType: esnats.Memory,
		}),
		snapnats.NewSnapshotStore[customer.AggregateRoot](ctx, js, snapnats.SnapshotStoreConfig{
			StoreType: snapnats.Memory,
		}),
		essrv.AggregateConfig{
			SnapthotMsgThreshold: 5,
		},

		essrv.WithEvent[customer.OrderRejected](),
		essrv.WithEvent[customer.OrderAccepted](),
		essrv.WithEvent[customer.Registered](),
	)

	ord := natsstore.NewAggregate(ctx, js,
		natsstore.NatsAggregateConfig{
			AggregateConfig: essrv.AggregateConfig{
				SnapthotMsgThreshold: 5,
			},
			StoreType: natsstore.Memory,
		},
		essrv.WithEvent[order.Closed](),
		essrv.WithEvent[order.Posted](),
		essrv.WithEvent[order.Verified](),
	)
	mod := &Module{
		PostOrder:          orderUsecase.NewPostOrderHandler(ord),
		RegisterCustomer:   customerUsecase.NewRegisterHandler(cust),
		OrderPostedHandler: customerUsecase.NewOrderPostedHandler(customerUsecase.NewVerifyOrderHandler(cust)),
		OrderStream:        ord,
	}

	return mod
}

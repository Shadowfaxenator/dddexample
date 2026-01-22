package sales

import (
	"context"
	"log/slog"
	"time"

	"github.com/alekseev-bro/ddd/pkg/aggregate"

	"github.com/alekseev-bro/ddd/pkg/natsstore"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	customercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/command"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	ordercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/command"

	"github.com/nats-io/nats.go/jetstream"
)

type Module struct {
	RegisterCustomer customercmd.CustomerRegisterHandler
	PostOrder        ordercmd.OrderpostHandler
	OrderStream      aggregate.Subscriber[order.Order]
	CustomerStream   aggregate.Subscriber[customer.Customer]
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {
	var cons []aggregate.Drainer
	cust := natsstore.NewStore(ctx, js,
		natsstore.WithInMemory[customer.Customer](),
		natsstore.WithSnapshot[customer.Customer](5, time.Second),
		natsstore.WithEvent[customer.OrderRejected, customer.Customer]("OrderRejected"),
		natsstore.WithEvent[customer.OrderAccepted, customer.Customer]("OrderAccepted"),
		natsstore.WithEvent[customer.Registered, customer.Customer]("CustomerRegistered"),
	)

	ord := natsstore.NewStore(ctx, js,
		natsstore.WithInMemory[order.Order](),
		natsstore.WithSnapshot[order.Order](5, time.Second),
		natsstore.WithEvent[order.Closed, order.Order]("OrderClosed"),
		natsstore.WithEvent[order.Posted, order.Order]("OrderPosted"),
		natsstore.WithEvent[order.Verified, order.Order]("OrderVerified"),
	)

	d, err := aggregate.Project(ctx, ord, customercmd.NewOrderPostedHandler(
		customercmd.NewVerifyOrderHandler(cust),
	))
	if err != nil {
		slog.Error("subscription create consumer", "error", err)
		panic(err)
	}
	cons = append(cons, d)

	d, err = aggregate.Project(ctx, cust, ordercmd.NewOrderRejectedHandler(
		ordercmd.NewCloseOrderHandler(ord),
	))
	if err != nil {
		slog.Error("subscription create consumer", "error", err)
		panic(err)
	}
	cons = append(cons, d)

	mod := &Module{
		PostOrder:        ordercmd.NewPostOrderHandler(ord),
		RegisterCustomer: customercmd.NewRegisterHandler(cust),
		OrderStream:      ord,
		CustomerStream:   cust,
	}

	go func() {
		<-ctx.Done()
		for _, c := range cons {
			if err := c.Drain(); err != nil {
				slog.Error("subscription drain", "error", err)
			}
		}
		slog.Info("all drainded")
	}()

	return mod
}

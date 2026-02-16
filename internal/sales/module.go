package sales

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/drivers/stream/esnats"
	"github.com/alekseev-bro/ddd/pkg/stream"

	"github.com/alekseev-bro/ddd/pkg/natsstore"

	"github.com/alekseev-bro/dddexample/contracts/v1/carpark"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	customercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/command"
	custquery "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/query"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	ordercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/command"
	orderquery "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/query"

	"github.com/nats-io/nats.go/jetstream"
)

type Projector interface {
	Project(any) error
}

type Module struct {
	RegisterCustomer aggregate.CommandHandler[customer.Customer, customercmd.Register]
	PostOrder        aggregate.CommandHandler[order.Order, ordercmd.Post]
	OrderStream      aggregate.Subscriber[order.Order]
	CustomerStream   aggregate.Subscriber[customer.Customer]
	OrderProjection  orderquery.OrdersLister
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {
	var cons []stream.Drainer
	cust, err := natsstore.New(ctx, js,
		natsstore.WithInMemory[customer.Customer](),
		natsstore.WithSnapshot[customer.Customer](5, time.Second, 5*time.Second),
		natsstore.WithEvent[customer.OrderRejected, customer.Customer]("OrderRejected"),
		natsstore.WithEvent[customer.OrderAccepted, customer.Customer]("OrderAccepted"),
		natsstore.WithEvent[customer.Registered, customer.Customer]("CustomerRegistered"),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ord, err := natsstore.New(ctx, js,
		natsstore.WithInMemory[order.Order](),
		natsstore.WithSnapshot[order.Order](5, time.Second, 5*time.Second),
		natsstore.WithEvent[order.Closed, order.Order]("OrderClosed"),
		natsstore.WithEvent[order.Posted, order.Order]("OrderPosted"),
		natsstore.WithEvent[order.Verified, order.Order]("OrderVerified"),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	d, err := aggregate.ProjectEvent(ctx, ord, customercmd.NewOrderPostedHandler(
		customercmd.NewVerifyOrderHandler(cust),
	))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cons = append(cons, d)

	d, err = aggregate.ProjectEvent(ctx, cust, ordercmd.NewOrderRejectedHandler(
		ordercmd.NewCloseOrderHandler(ord),
	))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cons = append(cons, d)
	custproj := custquery.NewCustomerProjection()
	ordproj := orderquery.NewMemOrders()
	d, err = ord.Subscribe(ctx, orderquery.NewOrderListProjector(custproj, ordproj))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cons = append(cons, d)

	d, err = cust.Subscribe(ctx, custquery.NewCustomerListProjector(custproj))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cons = append(cons, d)

	d, err = cust.Subscribe(ctx, custquery.NewCustomerListProjector(custproj))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cons = append(cons, d)
	dr, err := esnats.NewDriver(ctx, js, "car", esnats.WithStoreType(esnats.Memory))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	carStream := stream.New(ctx, dr, stream.WithEvent[carpark.CarArrived]("CarArrived"))
	_ = carStream
	// carStream.Subscribe(ctx, nil)

	mod := &Module{
		PostOrder:        ordercmd.NewPostOrderHandler(ord),
		RegisterCustomer: customercmd.NewRegisterHandler(cust),
		OrderStream:      ord,
		CustomerStream:   cust,
		OrderProjection:  ordproj,
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

package sales

import (
	"context"
	"log/slog"
	"os"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/drivers/stream/esnats"
	"github.com/alekseev-bro/ddd/pkg/stream"

	na "github.com/alekseev-bro/ddd/pkg/natsaggregate"

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
	var cons stream.DrainList
	cust, err := na.New(ctx, js,
		na.WithInMemory[customer.Customer](),
		na.WithSnapshotEventCount[customer.Customer](5),
		na.WithEvent[customer.OrderRejected, customer.Customer](),
		na.WithEvent[customer.OrderAccepted, customer.Customer](),
		na.WithEvent[customer.Registered, customer.Customer](),
	)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ord, err := na.New(ctx, js,
		na.WithInMemory[order.Order](),
		na.WithSnapshotEventCount[order.Order](5),
		na.WithEvent[order.Closed, order.Order](),
		na.WithEvent[order.Posted, order.Order](),
		na.WithEvent[order.Verified, order.Order](),
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
	dr, err := esnats.NewStore(ctx, js, "car", esnats.WithStoreType(esnats.Memory))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	carStream, err := stream.New(dr, stream.WithEvent[carpark.CarArrived]())
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
		if err := cons.Drain(); err != nil {
			slog.Error("subscription drain", "error", err)
			return
		}
		slog.Info("all drainded")
	}()

	return mod
}

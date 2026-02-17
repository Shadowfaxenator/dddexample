package sales

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/drivers/snapshot/natssnapshot"
	"github.com/alekseev-bro/ddd/pkg/drivers/stream/natsstream"
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
	RegisterCustomer aggregate.CommandHandler[customercmd.Register, customer.Customer]
	PostOrder        aggregate.CommandHandler[ordercmd.Post, order.Order]
	OrderStream      aggregate.Subscriber[order.Order]
	CustomerStream   aggregate.Subscriber[customer.Customer]
	OrderProjection  orderquery.OrdersLister
}

func NewModule(ctx context.Context, js jetstream.JetStream) *Module {
	cusstream, err := natsstream.New(ctx, js, "customer",
		natsstream.WithStoreType(natsstream.Memory))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	csnap, err := natssnapshot.New(ctx, js, "customer",
		natssnapshot.WithStoreType(natssnapshot.Memory))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cust, err := aggregate.New(ctx, cusstream, csnap,
		aggregate.WithSnapshotEventCount[customer.Customer](5),
		aggregate.WithEvent[customer.OrderRejected, customer.Customer](),
		aggregate.WithEvent[customer.OrderAccepted, customer.Customer](),
		aggregate.WithEvent[customer.Registered, customer.Customer](),
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
	if err := aggregate.ProjectEvent(ctx, ord, customercmd.NewOrderPostedHandler(
		customercmd.NewVerifyOrderHandler(cust),
	)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err = aggregate.ProjectEvent(ctx, cust, ordercmd.NewOrderRejectedHandler(
		ordercmd.NewCloseOrderHandler(ord),
	)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	custproj := custquery.NewCustomerProjection()
	ordproj := orderquery.NewMemOrders()

	if err = ord.Subscribe(ctx, orderquery.NewOrderListProjector(custproj, ordproj)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err = cust.Subscribe(ctx, custquery.NewCustomerListProjector(custproj)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err = cust.Subscribe(ctx, custquery.NewCustomerListProjector(custproj)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	es, err := natsstream.New(ctx, js, "car", natsstream.WithStoreType(natsstream.Memory))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	carStream, err := stream.New(es, stream.WithEvent[carpark.CarArrived]())
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
		wg := new(sync.WaitGroup)
		wg.Go(func() {
			cust.Drain()
		})
		wg.Go(func() {
			ord.Drain()
		})
		wg.Go(func() {
			carStream.Drain()
		})
		wg.Wait()
		slog.Info("all drainded")
	}()

	return mod
}

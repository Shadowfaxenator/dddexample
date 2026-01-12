package main

import (
	"context"
	"time"

	"os"
	"os/signal"

	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/dddexample/internal/sales"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	customercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/command"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	ordercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/command"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"

	"github.com/google/uuid"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	//	slog.SetLogLoggerLevel(slog.LevelWarn)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)

	if err != nil {
		panic(err)
	}

	s := sales.NewModule(ctx, js)
	events.Project(ctx, s.OrderPostedHandler)
	time.Sleep(time.Second)
	custid := values.CustomerID(uuid.New())
	cmdCust := customercmd.Register{
		ID:   custid,
		Name: "Joe",
		Age:  16,
	}

	err = s.RegisterCustomer.Handle(ctx, events.ID[customer.Customer](custid), cmdCust, custid.String())
	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, custid.String(), &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for range 1 {

		// ordid := s.Order.NewID()
		// idempo := aggregate.NewUniqueCommandIdempKey[*sales.CreateOrder](ordid)

		// _, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
		// if err != nil {
		// 	panic(err)
		// }
		ordID := values.OrderID(uuid.New())
		ordCmd := ordercmd.Post{
			ID:         ordID,
			CustomerID: custid,
		}
		err = s.PostOrder.Handle(ctx, events.ID[order.Order](ordID), ordCmd, ordID.String())
		if err != nil {
			panic(err)
		}
	}
	<-ctx.Done()

}

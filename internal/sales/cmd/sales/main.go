package main

import (
	"context"
	"time"

	"os"
	"os/signal"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/orders"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
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
	essrv.ProjectEvent(ctx, s.OrderPostedHandler)
	time.Sleep(time.Second)
	custid := domain.CustomerID(uuid.New())
	cust := &customers.Customer{ID: custid, Name: "Joe", Age: 21}
	err = s.RegisterCustomer.Handle(ctx, essrv.ID[customers.Customer](custid), customer.Register{Customer: cust}, custid.String())
	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, custid.String(), &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for range 5 {

		// ordid := s.Order.NewID()
		// idempo := aggregate.NewUniqueCommandIdempKey[*sales.CreateOrder](ordid)

		// _, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
		// if err != nil {
		// 	panic(err)
		// }
		ordID := domain.OrderID(uuid.New())
		ord := &orders.Order{ID: ordID, CustomerID: custid}
		err = s.PostOrder.Handle(ctx, essrv.ID[orders.Order](ordID), order.Post{Order: ord}, ordID.String())
		if err != nil {
			panic(err)
		}
	}
	<-ctx.Done()

}

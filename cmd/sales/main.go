package main

import (
	"context"

	"os"
	"os/signal"

	"github.com/alekseev-bro/ddd/pkg/domain"
	"github.com/alekseev-bro/dddexample/internal/domain/sales"

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

	s := sales.New(ctx, js)
	s.StartOrderCreationSaga(ctx)
	s.StartProjections(ctx)
	custid := s.Customer.NewID()
	_, err = s.Customer.Execute(ctx, custid.String(), &domain.Create[sales.Customer]{
		ID:   custid,
		Body: &sales.Customer{Name: "dddds", Age: 20},
	})
	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, idempc, &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for range 10 {

		// ordid := s.Order.NewID()
		// idempo := aggregate.NewUniqueCommandIdempKey[*sales.CreateOrder](ordid)

		// _, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
		// if err != nil {
		// 	panic(err)
		// }
		orid := s.Order.NewID()
		_, err := s.Order.Execute(ctx, orid.String(), &domain.Create[sales.Order]{
			ID:   orid,
			Body: &sales.Order{CustomerID: custid},
		})
		if err != nil {
			panic(err)
		}
	}

	<-ctx.Done()

}

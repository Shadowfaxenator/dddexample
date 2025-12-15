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

	cusid := s.Customer.NewID()

	idempc := domain.NewIdempotencyKey(cusid, "CreateCustomer")

	_, err = s.Customer.Execute(ctx, idempc, &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	if err != nil {
		panic(err)
	}
	for range 300 {

		ordid := s.Order.NewID()
		idempo := domain.NewIdempotencyKey(ordid, "CreateOrder")

		_, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
		if err != nil {
			panic(err)
		}
	}

	<-ctx.Done()

}

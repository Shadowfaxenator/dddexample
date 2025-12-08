package main

import (
	"context"
	"ddd/pkg/domain"

	"dddexample/internal/domain/sales"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	//	slog.SetLogLoggerLevel(slog.LevelError)
	// nc, err := nats.Connect(nats.DefaultURL)
	// if err != nil {
	// 	slog.Error("connect to nats", "error", err)
	// 	panic(err)
	// }

	// _, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: "atest", Subjects: []string{"atest.>"}, AllowAtomicPublish: true})
	// if err != nil {
	// 	slog.Error("create stream", "error", err)
	// 	panic(err)
	// }

	// _, err = js.PublishMsg(ctx, m, jetstream.WithExpectLastSequenceForSubject(uint64(0), "atest.t"))
	// if err != nil {
	// 	slog.Error("publish message", "error", err)
	// 	panic(err)
	// }

	// w.Start()

	s := sales.New(ctx)

	go func() {
		for {

			cusid := domain.NewID[sales.Customer]()
			idempc := domain.NewIdempotencyKey(cusid, "CreateCustomer")

			err := s.Customer.Execute(ctx, idempc, &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
			if err != nil {
				panic(err)
			}
			// for range 1 {

			ordid := domain.NewID[sales.Order]()
			idempo := domain.NewIdempotencyKey(ordid, "CreateOrder")

			err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
			if err != nil {
				panic(err)
			}
			<-time.After(1 * time.Second)
		}

	}()

	<-ctx.Done()

}

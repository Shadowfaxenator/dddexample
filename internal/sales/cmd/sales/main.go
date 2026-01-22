package main

import (
	"context"
	"log/slog"
	"time"

	"os"
	"os/signal"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	customercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/command"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	ordercmd "github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order/command"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	slog.SetLogLoggerLevel(slog.LevelInfo)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	js, err := jetstream.New(nc)

	if err != nil {
		panic(err)
	}

	s := sales.NewModule(ctx, js)
	time.Sleep(time.Second)
	c := customer.New("Joe", 16, nil)
	cmdCust := customercmd.Register{Customer: c}
	ctxIdemp := aggregate.ContextWithIdempotancyKey(ctx, c.ID.String())
	_, err = s.RegisterCustomer.Handle(ctxIdemp, cmdCust)
	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, custid.String(), &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for {
		select {
		case <-ctx.Done():
			time.Sleep(time.Second * 3)
			return

		case <-time.After(time.Second * 2):
			// ordid := s.Order.NewID()
			// idempo := aggregate.NewUniqueCommandIdempKey[*sales.CreateOrder](ordid)

			// _, err = s.Order.Execute(ctx, idempo, &sales.CreateOrder{OrderID: ordid, CustID: cusid})
			// if err != nil {
			// 	panic(err)
			// }
			o := order.New(c.ID, nil)
			ordCmd := ordercmd.Post{
				Order: o,
			}
			ctxIdemp := aggregate.ContextWithIdempotancyKey(ctx, o.ID.String())

			_, err = s.PostOrder.Handle(ctxIdemp, ordCmd)
			if err != nil {
				panic(err)
			}
		}
	}

	//	<-ctx.Done()

}

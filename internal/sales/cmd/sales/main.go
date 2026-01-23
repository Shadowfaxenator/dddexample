package main

import (
	"context"
	"fmt"
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
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	slog.SetLogLoggerLevel(slog.LevelError)

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
	_, err = s.RegisterCustomer.HandleCommand(ctxIdemp, cmdCust)
	if err != nil {
		panic(err)
	}

	//	idempc := aggregate.NewUniqueCommandIdempKey[*sales.CreateCustomer](cusid)

	// _, err = s.Customer.Execute(ctx, custid.String(), &sales.CreateCustomer{Customer: sales.Customer{ID: cusid, Name: "John", Age: 20}})
	// if err != nil {
	// 	panic(err)
	// }
	for range 10 {

		o := order.New(c.ID, order.OrderLines{
			order.OrderLine{
				CarID:    aggregate.NewID(),
				Quantity: 1,
				Price:    values.NewMoney("USD", 200, 2),
			},
			order.OrderLine{
				CarID:    aggregate.NewID(),
				Quantity: 2,
				Price:    values.NewMoney("USD", 100, 2),
			},
		})
		ordCmd := ordercmd.Post{
			Order: o,
		}
		ctxIdemp := aggregate.ContextWithIdempotancyKey(ctx, o.ID.String())

		_, err = s.PostOrder.HandleCommand(ctxIdemp, ordCmd)
		if err != nil {
			panic(err)
		}
	}
	<-time.After(time.Second * 2)
	l, err := s.OrderProjection.ListAll()
	if err != nil {
		panic(err)
	}
	for _, o := range l {
		fmt.Printf("order: %+v\n", o)
	}

	<-ctx.Done()

}

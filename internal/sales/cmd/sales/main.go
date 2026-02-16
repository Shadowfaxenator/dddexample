package main

import (
	"context"
	"fmt"
	"log/slog"
	"syscall"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	slog.SetLogLoggerLevel(slog.LevelInfo)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()
	js, err := jetstream.New(nc)

	if err != nil {
		panic(err)
	}

	s := sales.NewModule(ctx, js)
	time.Sleep(time.Second)
	c := customer.New("Joe", 16, nil)
	cmdCust := customercmd.Register{Customer: c}
	_, err = s.RegisterCustomer.HandleCommand(ctx, cmdCust)
	if err != nil {
		panic(err)
	}

	for range 10 {
		carid1, err := aggregate.NewID()
		if err != nil {
			panic(err)
		}
		carid2, err := aggregate.NewID()
		if err != nil {
			panic(err)
		}
		o := order.New(c.ID, order.OrderLines{
			order.OrderLine{
				CarID:    carid1,
				Quantity: 1,
				Price:    values.NewMoney("USD", 200, 2),
			},
			order.OrderLine{
				CarID:    carid2,
				Quantity: 2,
				Price:    values.NewMoney("USD", 100, 2),
			},
		})
		ordCmd := ordercmd.Post{
			Order: o,
		}

		_, err = s.PostOrder.HandleCommand(ctx, ordCmd)
		if err != nil {
			panic(err)
		}
	}
	<-time.After(time.Second * 2)
	l, err := s.OrderProjection.ListOrders()
	if err != nil {
		panic(err)
	}
	fmt.Println("Projection Example (OrderList in memory) :")
	for _, o := range l {
		fmt.Printf("order: %+v\n", o)
	}

	<-ctx.Done()

}

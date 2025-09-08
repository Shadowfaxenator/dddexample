package main

import (
	"context"
	"ddd/pkg/aggregate"
	"encoding/json"
	"os"
	"os/signal"
	"ttt/internal/domain/sales"
)

type My struct {
	Name string
}

type Applyer[T any] interface {
	Apply(t *T)
}
type Event[T any] struct {
	Body Applyer[T]
	Type string
}
type RawEvent struct {
	Body json.RawMessage
	Type string
}

type MyEvent[T any] struct {
	Name string
	Cars []string
}

func (e MyEvent[T]) Apply(t *My) {
	t.Name = e.Name

}

type Car struct {
	ID    string
	Color string
	User  string
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	s := sales.NewSubDomain(ctx)
	cusid := aggregate.NewID()
	err := s.CreateCustomer(ctx, cusid, "Bob", 22)
	if err != nil {
		panic(err)
	}
	ordid := aggregate.NewID()
	err = s.CreateOrder(ctx, ordid, cusid)
	if err != nil {
		panic(err)
	}

	//err = s.CloseOrder(ctx, ordid)
	// if err != nil {
	// 	panic(err)
	// }

	<-ctx.Done()
}

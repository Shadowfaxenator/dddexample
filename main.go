package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
	"ttt/gonvex"
	"ttt/internal/domains/sales"
)

type My struct {
	Name string
	As   []sales.Address
	A    sales.Address
	Can  bool
}

func main() {
	//my := My{Name: "John", As: []sales.Address{{Street: "123 Main St", City: "Anytown", State: "CA", Zip: "12345"}}, Address: sales.Address{Street: "456 Oak St", City: "Othertown", State: "TX", Zip: "67890"}}
	// connect

	//a := um["Can"].(string)
	//fmt.Printf("um: %v\n", a)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	customer := sales.NewCustomerAggregate(ctx)
	id, _ := customer.Create(ctx, "Artem")
	fmt.Printf("id: %v\n", id)
	//id, _ := gonvex.AggregateIDFromString("fbfcb781-514b-48c2-b868-50ffacee1dec")
	go func() {
		for {
			select {
			case <-time.After(3 * time.Second):
				err := customer.AddCar(ctx, id)
				if err != nil {
					fmt.Printf("error adding car: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	err := customer.StartProjection(ctx, "customer_projection", func(ev gonvex.Event) any {
		switch e := ev.(type) {
		case *sales.CustomerCreated:
			fmt.Printf("customer created: %v\n", e)
		case *sales.CustomerCarAdded:
			fmt.Printf("car added to customer: %v\n", e)
		}
		return ""
	})
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
}

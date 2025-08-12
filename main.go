package main

import (
	"context"
	"ttt/gonvex"
	"ttt/internal/domains/sales"
)

func main() {
	ctx := context.Background()
	customer := sales.NewCustomerAggregate(ctx)
	//agr.Root = NewCustomer("Joe")
	event := customer.Create.New(sales.CustomerCreated{Name: "Joe", Cars: make([]gonvex.AggregateID, 0)})
	// if len(a.cars) >= 10 {
	// 	return errors.New("not allowed to add more than 10 cars")
	// }
	customer.Stream.NewAggregate().SendEvent(ctx, event)
	// custid, err := agr.Create(ctx)
	// if err != nil {
	// 	panic(err)

	//	fmt.Printf("out: %v\n", ev)
	// for _, v := range e {
	// 	fmt.Printf("reflect.TypeOf(ev).Elem().Name(): %v\n", reflect.TypeOf(v).Name())
	// }

	// //ee, _ := json.Marshal(&My{Name: "nnnn"})
	// var ev = &Event{ID: uuid.New().String(), Payload: &My{Name: "nnnn"}}
	// b, _ := json.Marshal(ev)
	// var event Event
	// json.Unmarshal(b, &event)
	// fmt.Printf("event: %+v\n", event)
	// var gg any
	// //json.Unmarshal(event.Payload, &gg)

	// fmt.Printf("event: %+v\n", gg)

	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// defer cancel()

	// perid, _ := gonvex.AggregateIDFromString("6c3a7819-e512-4e4d-824c-7e00aaaf5f77")
	// carid, _ := gonvex.AggregateIDFromString("82de484a-a864-4d2f-801b-805e9920d0bc")

	// hr := sales.New(ctx)
	// if err := hr.AddCarToPerson(ctx, perid, carid); err != nil {
	// 	log.Println(err)

	// }
	// <-ctx.Done()
	// fmt.Println(ctx.Err())
	// select {}
}

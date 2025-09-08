package aggregate

import (
	"context"
	"fmt"
)

// func SubscribeAll[T Reducible[T], U AggregateRoot[T]](ctx context.Context, name string, a U, handler func(event T) error, ordered bool) {

// }

func (a *Aggregate[T]) Subscribe(ctx context.Context, name string, handler func(Applyer[T]) error, ordered bool) {
	a.stream.Subscribe(ctx, name, func(b []byte) error {
		var rec StoreRecord
		a.serder.Deserialize(b, &rec)
		ev, err := a.eventRegistry.Get(rec.Type, rec.Body)
		if err != nil {
			return fmt.Errorf("subscribe: %w", err)
		}

		return handler(ev.(Applyer[T]))

	}, ordered)
	//t := reflect.TypeFor[T]()

	//fmt.Printf("ParseToSql(t): %+v\n", ParseToSql(t))

}

package gonvex

type AggregateType string

// type BoundedContext[T []Store[T]] struct {
// 	Name   string
// 	Stores []map[AggregateType]T
// }

// func NewBoundedContext[T []Store[Reducible]](name string, stores ...Store[Reducible]) {
// 	bc := &BoundedContext[T]{Name: name, Stores: make([]map[AggregateType]T, 0)}
// 	for _, v := range stores {

// 		bc.Stores[AggregateType(v.Type())] = T[]
// 		bc.Stores[AggregateType(v.Type())].SetBC(name)
// 	}
// 	//return bc
// }

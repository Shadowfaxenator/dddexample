package ddd

import (
	"context"
	"ddd/pkg/aggregate"
)

type Root[T any] interface {
	Commander[T]
	Subscriber[T]
}

type Commander[T any] interface {
	Command(context.Context, aggregate.ID, aggregate.Executer[T]) error
	CommandFunc(context.Context, aggregate.ID, func(*T) (*aggregate.Event[T], error)) error
}
type Subscriber[T any] interface {
	Subscribe(ctx context.Context, name string, handler func(aggregate.Applyer[T]) error, ordered bool)
}

func NewAggregate[T any](ctx context.Context) *aggregate.Aggregate[T] {
	return aggregate.New[T](ctx)
}

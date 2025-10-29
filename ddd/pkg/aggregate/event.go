package aggregate

// EventError is not saved to the event store.
type EventError[T any] struct {
	Reason string
}

func (e EventError[T]) Error() string {
	return e.Reason
}

func (e EventError[T]) Apply(*T) {}

type EventRegistry[T any] interface {
	RegisterEvent(Event[T])
}

type Event[T any] interface {
	Apply(*T)
}

func RegisterEvent[E Event[T], T any](reg EventRegistry[T]) {
	var ev E
	reg.RegisterEvent(ev)
}

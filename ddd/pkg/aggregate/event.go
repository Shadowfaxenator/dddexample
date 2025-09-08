package aggregate

import (
	"reflect"
)

func typeFromName(e any) string {
	t := reflect.TypeOf(e)
	switch t.Kind() {
	case reflect.Pointer:
		return t.Elem().Name()
	default:
		return t.Name()
		//	json.Marshal()
	}
}

type Event[T any] struct {
	Applyer[T]
	Type string
}

func NewEvent[T any](event Applyer[T]) *Event[T] {

	return &Event[T]{Applyer: event, Type: typeFromName(event)}
}

type EventRegistry[T any] interface {
	RegisterEvent(Applyer[T])
}

type Applyer[T any] interface {
	Apply(*T)
}

func RegisterEvent[E Applyer[T], T any](reg EventRegistry[T]) {
	var ev E
	reg.RegisterEvent(ev)
}

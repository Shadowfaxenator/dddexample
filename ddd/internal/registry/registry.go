package registry

import (
	"ddd/internal/serde"
	"fmt"
	"log/slog"

	"reflect"
)

type ctor[T any] func(payload []byte) (T, error)

type RegistryStore[T any] struct {
	serder serde.Serder
	items  map[string]ctor[T]
}

func (r *RegistryStore[T]) Register(item T) {
	t := reflect.TypeOf(item)
	slog.Info("EventRegistered", "type", t.Name())

	if t.Kind() != reflect.Struct && t.Kind() != reflect.Interface {
		panic("register: registered type must be struct or interface")
	}
	ctor := func(payload []byte) (T, error) {
		var zero T
		vt := reflect.New(t).Interface()
		if err := r.serder.Deserialize(payload, vt); err != nil {
			return zero, fmt.Errorf("registry: %w", err)
		}
		return vt.(T), nil
		//	var val V

	}

	r.items[TypeNameFrom(item)] = ctor
}

func TypeNameFrom(e any) string {
	if strev, ok := e.(fmt.Stringer); ok {
		return strev.String()
	}

	t := reflect.TypeOf(e)
	switch t.Kind() {
	case reflect.Struct:
		return t.Name()
	case reflect.Pointer:
		return t.Elem().Name()
	default:
		panic("unsupported type")

		//	json.Marshal()
	}
}
func (r *RegistryStore[T]) Exists(tname string) bool {
	if _, ok := r.items[tname]; ok {
		return true
	}
	return false
}

func (r *RegistryStore[T]) Get(tname string, b []byte) (T, error) {

	if ct, ok := r.items[tname]; ok {

		tt, err := ct(b)
		if err != nil {
			return tt, fmt.Errorf("registry: %w", err)
		}
		return tt, nil
	}
	slog.Error("registry: no type found", "type", tname)
	panic("unrecovered")
	//return nil, fmt.Errorf(")
}

func New[T any](s serde.Serder) *RegistryStore[T] {
	return &RegistryStore[T]{items: make(map[string]ctor[T]), serder: s}
}

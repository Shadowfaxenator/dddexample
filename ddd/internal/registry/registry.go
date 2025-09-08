package registry

import (
	"ddd/internal/serde"
	"fmt"
	"log/slog"

	"reflect"
)

type ctor func(payload []byte) (any, error)

type RegistryStore struct {
	serder serde.Serder
	items  map[string]ctor
}

func (r *RegistryStore) Register(item any) {
	t := reflect.TypeOf(item)
	slog.Info("EventRegistered", "type", t.Name())

	if t.Kind() != reflect.Struct && t.Kind() != reflect.Interface {
		panic("register: registered type must be struct or interface")
	}
	ctor := func(payload []byte) (any, error) {

		vt := reflect.New(t).Interface()
		if err := r.serder.Deserialize(payload, vt); err != nil {
			return nil, fmt.Errorf("registry: %w", err)
		}
		return vt, nil
		//	var val V

	}
	r.items[t.Name()] = ctor
}

func (r *RegistryStore) Get(tname string, b []byte) (any, error) {

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

func New(s serde.Serder) *RegistryStore {
	return &RegistryStore{items: make(map[string]ctor), serder: s}
}

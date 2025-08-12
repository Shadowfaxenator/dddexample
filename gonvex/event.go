package gonvex

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

const EVENT_HEADER = "event_type"

type ctor func(payload []byte) (any, error)

var EventRegistry = make(map[string]ctor)

type Event[T Reducible] interface {
	Type() string
	Data() (json.RawMessage, error)
	Invariant() Invariant[T]
}

type AggregateEvent struct {
	ID       uuid.UUID
	BContext string    `json:"bounded_context"`
	AggrID   uuid.UUID `json:"aggregate_id"`
	AggrType string    `json:"aggregate_type"`
	Event    json.RawMessage
}

type CoreEvent[T any, U Reducible] struct {
	EType   string `json:"event_type"`
	Payload T      `json:"payload"`
	inv     Invariant[U]
}

func (e *CoreEvent[T, U]) Type() string {
	return e.EType
}

func (e *CoreEvent[T, U]) Data() (json.RawMessage, error) {

	return json.Marshal(e.Payload)
}

func (e *CoreEvent[T, U]) Invariant() Invariant[U] {

	return e.inv
}

type EventType[T any, U Reducible] struct {
	etype string
	inv   Invariant[U]
}

func (e EventType[T, U]) New(args T) Event[U] {
	event := CoreEvent[T, U]{EType: e.etype, Payload: args, inv: e.inv}

	return &event
}

type Command[T any, U Reducible] interface {
	New(args T) Event[U]
}
type Invariant[T Reducible] func(aggr T) error

func NewCommand[T any, U Reducible](inv Invariant[U]) Command[T, U] {
	name := reflect.TypeFor[T]().Name()
	EventRegistry[name] = func(payload []byte) (any, error) {
		ee, err := Deserialize[AggregateEvent](payload)
		if err != nil {
			return nil, fmt.Errorf("new event factory: %w", err)
		}
		des, err := Deserialize[T](ee.Event)
		if err != nil {
			return nil, fmt.Errorf("new event factory: %w", err)
		}
		return *des, nil
	}
	ev := EventType[T, U]{etype: name, inv: inv}
	return &ev
}

func Deserialize[T any](b []byte) (*T, error) {
	var ev T
	if err := json.Unmarshal(b, &ev); err != nil {
		return nil, fmt.Errorf("deserialize: %w", err)
	}
	return &ev, nil
}

func Serialize(e any) ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("serialize: %w", err)
	}
	return b, nil
}

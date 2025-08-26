package gonvex

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

const (
	EVENT_TYPE_HEADER        = "ev_type"
	AGGREGATE_ID_HEADER      = "agg_id"
	MAX_MSGS            uint = 100
)

type ctor func(payload []byte) (Event, error)

var EventRegistry = make(map[string]ctor)

type Event interface {
	ID() ID
	Type() string
}

func registerEvent(event Event, bc string) {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Pointer {
		panic("registered event can't be a pointer")
	}
	vt := reflect.New(t).Interface()

	EventRegistry[eventName(event.Type(), bc)] = func(payload []byte) (Event, error) {
		var ae AggregateEvent
		err := Deserialize(payload, &ae)
		if err != nil {
			return nil, fmt.Errorf("new event factory: %w", err)
		}

		if err := Deserialize(ae.Event, vt); err != nil {

			return nil, fmt.Errorf("new event factory: %w", err)
		}

		return vt.(Event), nil
	}
}

type AggregateEvent struct {
	ID       uuid.UUID
	BContext string    `json:"bounded_context"`
	AggrID   uuid.UUID `json:"aggregate_id"`
	AggrType string    `json:"aggregate_type"`
	Event    []byte
}

//	type Command[Evt any, AgrRoot Reducible] interface {
//		Run(ctx context.Context, event Evt, aggr *Aggregate[AgrRoot]) error
//	}

func Deserialize(b []byte, out any) error {

	if err := json.Unmarshal(b, &out); err != nil {
		return fmt.Errorf("deserialize: %w", err)
	}

	return nil
}

func Serialize(e any) ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("serialize: %w", err)
	}
	return b, nil
}

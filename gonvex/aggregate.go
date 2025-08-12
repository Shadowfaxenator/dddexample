package gonvex

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type AggregateID = uuid.UUID

func NewAggregateID() AggregateID {
	return AggregateID(uuid.New())
}

func AggregateIDFromString(id string) (AggregateID, error) {
	aid, err := uuid.Parse(id)
	return AggregateID(aid), err
}

type Reducible interface {
	Reduce(event any) error
}

type Aggregate[T Reducible] struct {
	ID         AggregateID
	Version    uint64
	Type       string
	BoundedCtx string
	Root       T
}

// func (a *AggregateRoot) Aggregate() *AggregateRoot {
// 	return a
// }

//func (a *Aggregate[T]) RunCommand()

func (a *Aggregate[T]) SendEvent(ctx context.Context, event Event[T]) error {
	inv := event.Invariant()
	if err := inv(a.Root); err != nil {
		return fmt.Errorf("send event: %w", err)
	}
	data, err := event.Data()

	if err != nil {
		return fmt.Errorf("send event: %w", err)
	}
	ae := &AggregateEvent{
		ID:       uuid.New(),
		BContext: a.BoundedCtx,
		AggrID:   a.ID,
		AggrType: a.Type,
		Event:    data,
	}

	msg := nats.NewMsg(fmt.Sprintf("%s-%s.%s", a.Type, a.BoundedCtx, a.ID.String()))
	msg.Header.Add(EVENT_HEADER, event.Type())

	p, err := Serialize(ae)
	if err != nil {
		return fmt.Errorf("send event func: %w", err)
	}

	msg.Data = p

	_, err = js.PublishMsg(ctx, msg, jetstream.WithExpectLastSequencePerSubject(a.Version))
	if err != nil {
		return fmt.Errorf("send event func: %w", err)
	}
	return nil

}

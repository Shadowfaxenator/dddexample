package gonvex

import (
	"github.com/google/uuid"
)

type ID [16]byte

func NewID() ID {
	a := uuid.New()
	return ID(a)
}
func (id ID) String() string {
	return uuid.UUID(id).String()
}

func IDFromString(id string) (ID, error) {
	aid, err := uuid.Parse(id)
	return ID(aid), err
}

type Reducible interface {
	Reduce(event Event) error
	Events() []Event
}

type Aggregate[T Reducible] struct {
	Version uint64
	Data    T
}

type Versioned interface {
	Version() uint64
}

// func (a *AggregateRoot) Aggregate() *AggregateRoot {
// 	return a
// }

//func (a *Aggregate[T]) RunCommand()

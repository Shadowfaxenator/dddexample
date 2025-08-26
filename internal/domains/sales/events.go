package sales

import "ttt/gonvex"

type CustomerCarAdded struct {
	AggID gonvex.ID
	CarID gonvex.ID
}

func (CustomerCarAdded) Type() string {
	return "CustomerCarAdded"
}
func (c CustomerCarAdded) ID() gonvex.ID {
	return c.AggID
}

type CustomerCreated struct {
	AggID    gonvex.ID
	Customer Customer
}

func (CustomerCreated) Type() string {
	return "CustomerCreated"
}

func (c CustomerCreated) ID() gonvex.ID {
	return c.AggID
}

// var (
// 	PersonCreated = func(person *Person) *gonvex.Event[any] {
// 		return gonvex.NewEvent("PERSON_CREATED", BOUNDED_CONTEXT, person)
// 	}
// 	PersonalCarAdded = func(carid uuid.UUID) *gonvex.Event[uuid.UUID] {
// 		return gonvex.NewEvent("PERSON_CAR_ADDED", BOUNDED_CONTEXT, carid)
// 	}
// 	PersonalCarRemoved = func(carid uuid.UUID) *gonvex.Event[uuid.UUID] {
// 		return gonvex.NewEvent("PERSON_CAR_REMOVED", BOUNDED_CONTEXT, carid)
// 	}
// )

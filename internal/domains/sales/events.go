package sales

import "ttt/gonvex"

const (
	BOUNDED_CONTEXT       = "PERSONAL"
	CUSTOMER_CREATED      = "CUSTOMER_CREATED"
	CUSTOMER_CAR_ADDED    = "CUSTOMER__CAR_ADDED"
	CUSTOMER__CAR_REMOVED = "CUSTOMER__CAR_REMOVED"
	ORDER_CREATED         = "ORDER_CREATED"
)

type CustomerCarAdded gonvex.AggregateID

type CustomerCreated Customer

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

package sales

import (
	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

type RentOrderStatus uint8

const (
	ProcessingByCustomer RentOrderStatus = iota
	ValidForProcessing
	Closed
)

type Order struct {
	CustomerID eventstore.ID[Customer]
	Cars       map[eventstore.ID[Car]]struct{}
	Status     RentOrderStatus
	Deleted    bool
}

package sales

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type RentOrderStatus uint8

const (
	ProcessingByCustomer RentOrderStatus = iota
	ValidForProcessing
	Closed
)

type Order struct {
	CustomerID aggregate.ID[Customer]
	Cars       map[aggregate.ID[Car]]struct{}
	Status     RentOrderStatus
}

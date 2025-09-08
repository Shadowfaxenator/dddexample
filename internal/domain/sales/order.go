package sales

import (
	"ddd/pkg/aggregate"
)

type RentOrderStatus uint8

const (
	ProcessingByCustomer RentOrderStatus = iota
	ValidForProcessing
	Closed
)

type Order struct {
	aggregate.ID
	CustomerID aggregate.ID
	Cars       map[aggregate.ID]struct{}
	Status     RentOrderStatus
}

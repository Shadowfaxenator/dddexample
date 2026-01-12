package values

import (
	"github.com/alekseev-bro/ddd/pkg/events"
)

type customer struct{}
type order struct{}
type product struct{}
type car struct{}

type CustomerID = events.ID[customer]
type OrderID = events.ID[order]
type ProductID = events.ID[product]
type CarID = events.ID[car]

package domain

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
)

type customer struct{}
type order struct{}
type product struct{}
type car struct{}

type CustomerID = essrv.ID[customer]
type OrderID = essrv.ID[order]
type ProductID = essrv.ID[product]
type CarID = essrv.ID[car]

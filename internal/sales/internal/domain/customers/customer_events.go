package customers

import "github.com/alekseev-bro/dddexample/internal/sales/internal/domain"

type CustomerRegistered struct {
	*Customer
}

func (e CustomerRegistered) Evolve(c *Customer) {
	*c = *e.Customer
}

type CustomerOrderClosed struct {
	CustomerID domain.CustomerID
	OrderID    domain.OrderID
}

func (CustomerOrderClosed) Evolve(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	CustomerID domain.CustomerID
	OrderID    domain.OrderID
}

func (OrderAccepted) Evolve(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	OrderID domain.OrderID
	Reason  string
}

func (OrderRejected) Evolve(c *Customer) {

}

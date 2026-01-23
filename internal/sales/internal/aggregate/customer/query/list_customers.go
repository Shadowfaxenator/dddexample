package query

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type Customer struct {
	ID   aggregate.ID
	Name string
}

type CustomerProjection struct {
	Customers map[aggregate.ID]Customer
}

func NewCustomerProjection() *CustomerProjection {
	return &CustomerProjection{
		Customers: make(map[aggregate.ID]Customer),
	}
}

func (store *CustomerProjection) GetCustomer(id aggregate.ID) (*Customer, bool) {
	customer, ok := store.Customers[id]
	return &customer, ok
}

func (store *CustomerProjection) AddCustomer(customer Customer) {
	store.Customers[customer.ID] = customer
}

func (store *CustomerProjection) UpdateCustomer(customer Customer) {
	store.Customers[customer.ID] = customer
}

func (store *CustomerProjection) DeleteCustomer(id aggregate.ID) {
	delete(store.Customers, id)
}

func (store *CustomerProjection) ListCustomers() []Customer {
	customers := make([]Customer, 0, len(store.Customers))
	for _, customer := range store.Customers {
		customers = append(customers, customer)
	}
	return customers
}

type CustomerListProjector struct {
	store *CustomerProjection
}

func NewCustomerListProjector(store *CustomerProjection) *CustomerListProjector {
	return &CustomerListProjector{
		store: store,
	}
}

func (h *CustomerListProjector) HandleEvents(ctx context.Context, event aggregate.Evolver[customer.Customer]) error {
	switch ev := event.(type) {
	case *customer.Registered:
		h.store.AddCustomer(Customer{
			ID:   ev.CustomerID,
			Name: ev.Name,
		})
	}
	return nil
}

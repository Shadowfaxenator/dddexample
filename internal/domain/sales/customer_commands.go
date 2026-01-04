package sales

import (
	estore "github.com/alekseev-bro/ddd/pkg/eventstore"
)

func (c *Customer) AddOrder() error {
	return nil
}

func (c *Customer) Create(name string, age uint) (CustomerEvents, error) {
	c.Name = name
	c.Age = age
	return estore.NewEvents(CustomerCreated{Customer: *c}), nil
}

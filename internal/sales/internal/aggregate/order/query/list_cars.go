package query

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/contracts/v1/carpark"
)

type CarsProjStore struct {
	cars []*Car
}

func NewCarsProjStore() *CarsProjStore {
	return &CarsProjStore{
		cars: make([]*Car, 0),
	}
}

func (s *CarsProjStore) AddCar(car *Car) error {
	s.cars = append(s.cars, car)
	return nil
}

func (s *CarsProjStore) ListCars() ([]*Car, error) {
	return s.cars, nil
}

type Car struct {
	ID  aggregate.ID
	VIN string
	CarModel
	RentState
	MaintananceState
}

type CarModel struct {
	Brand string
	Model string
}

type RentState uint

const (
	NotAvailable RentState = iota
	Available
)

type MaintananceState uint

const (
	InMaintenance MaintananceState = iota
	NotNeeded
	Needed
)

type carsLister interface {
	ListCars() ([]*Car, error)
}

type carAdder interface {
	AddCar(car *Car) error
}

func NewCarListProjector(store carAdder) *carListProjector {
	return &carListProjector{
		store: store,
	}
}

type carListProjector struct {
	store carAdder
}

func (p *carListProjector) Project(event any) error {
	switch ev := event.(type) {
	case *carpark.CarArrived:
		car := &Car{
			ID:  ev.Car.ID,
			VIN: ev.Car.VIN,
			CarModel: CarModel{
				Brand: ev.Car.Brand,
				Model: ev.Car.Model,
			},
			RentState:        RentState(ev.Car.RentState),
			MaintananceState: MaintananceState(ev.Car.MaintananceState),
		}
		return p.store.AddCar(car)
	default:
		return nil
	}
}

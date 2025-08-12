package carpark

// import (
// 	"context"
// 	"ttt/gonvex"
// )

// type BoundedContext struct {
// 	Cars gonvex.EventStream[*Car]
// 	//Stock gonvex.Store[*Stock]
// }

// func New(ctx context.Context) *BoundedContext {

// 	return &BoundedContext{
// 		Cars: gonvex.NewEventStream[Car](ctx, "Car", "CarPark"),
// 	}
// }

// func (h *BoundedContext) CreateCar(ctx context.Context, brand string, model string) (gonvex.AggregateID, error) {

// 	ca := h.Cars.New()
// 	// if len(a.cars) >= 10 {
// 	// 	return errors.New("not allowed to add more than 10 cars")
// 	// }

// 	if err := ca.SendEvent(ctx, CAR_CREATED, newCar(model, brand)); err != nil {
// 		return ca.ID, err
// 	}
// 	return ca.ID, nil
// 	//return agr.Mutate(ctx, gonvex.NewEvent(PERSON_CREATED, BOUNDED_CONTEXT, person))

// }

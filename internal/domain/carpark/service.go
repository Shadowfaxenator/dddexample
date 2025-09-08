package carpark

import (
	"context"
	"ddd/pkg/aggregate"
)

type Service struct {
	car *aggregate.Aggregate[Car]
}

func NewService(ctx context.Context) *Service {
	return &Service{
		car: aggregate.New[Car](ctx),
	}
}

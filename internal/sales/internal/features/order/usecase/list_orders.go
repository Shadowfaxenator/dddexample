package usecase

import (
	"time"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/money"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
)

type OrderListProjection struct {
	ID        ids.OrderID
	Total     money.Money
	CreatedAt time.Time
	UserName  string
}

type OrderListEventHandler struct {
	eventStore essrv.Root[order.AggregateRoot]
}

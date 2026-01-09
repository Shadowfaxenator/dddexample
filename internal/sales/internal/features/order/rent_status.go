package order

type RentOrderStatus uint8

const (
	StatusProcessingByCustomer RentOrderStatus = iota
	StatusValidForProcessing
	StatusClosed
)

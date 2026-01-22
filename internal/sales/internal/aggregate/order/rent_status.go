package order

type RentOrderStatus uint8

const (
	StatusNew RentOrderStatus = iota
	StatusProcessingByCustomer
	StatusValidForProcessing
	StatusClosed
)

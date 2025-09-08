package events

type SalesEvent uint

const (
	CustomerCreated SalesEvent = iota
	OrderCreated
)

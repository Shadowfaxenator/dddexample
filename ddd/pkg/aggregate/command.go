package aggregate

type Executer[T any] interface {
	Execute(*T) (*Event[T], error)
}

type CommandRegistry[T any] interface {
	RegisterCommand(Executer[T])
}

type Command[T any] struct {
	Executer[T]
	Type string
}

func NewCommand[T any](command Executer[T]) *Command[T] {
	return &Command[T]{Executer: command, Type: typeFromName(command)}
}

func RegisterCommand[E Executer[T], T any](root CommandRegistry[T]) {
	var com E
	root.RegisterCommand(com)
}

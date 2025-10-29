package aggregate

type Command[T any] interface {
	Execute(entity *T) Event[T]
}

type CommandRegistry[T any] interface {
	RegisterCommand(Command[T])
}

func RegisterCommand[E Command[T], T any](root CommandRegistry[T]) {
	var com E
	root.RegisterCommand(com)
}

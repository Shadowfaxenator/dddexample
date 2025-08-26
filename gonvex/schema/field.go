package schema

type Fildable interface {
	string | int | bool | float64 | float32 | []byte
}

type Field[T any] interface {
	Get() T
	Set(v T)
}

type List[T any] interface {
	Push(...T)
	Get() []T
}

func NewField[T any](value T) Field[T] {
	return &AggregateField[T]{Value: value}
}

func NewList[T any]() List[T] {
	return &AggregateList[T]{Value: make([]T, 0)}
}

type AggregateField[T any] struct {
	Value T
}

func (f *AggregateField[T]) Get() T {
	return f.Value
}

func (f *AggregateField[T]) Set(v T) {
	f.Value = v
}

type AggregateList[T any] struct {
	Value []T
}

func (f *AggregateList[T]) Push(v ...T) {
	f.Value = append(f.Value, v...)
}

func (f *AggregateList[T]) Get() []T {
	return f.Value
}

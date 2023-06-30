package minilang

// Maybe contains either value or an error.
type Maybe[T any] struct {
	Error error
	Value T
}

// NewError constructs a [Maybe] with the given error.
func NewError[T any](err error) Maybe[T] {
	return Maybe[T]{Error: err, Value: *new(T)}
}

// NewValue constructs a [Maybe] with the given value.
func NewValue[T any](value T) Maybe[T] {
	return Maybe[T]{Error: nil, Value: value}
}

package minidsl

// Maybe contains a value or an error. When processing the [Maybe], keep in mind that some
// errors, such as [ErrSkip] and [ErrException], have special meaning.
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

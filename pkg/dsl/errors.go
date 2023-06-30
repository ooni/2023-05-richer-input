package dsl

import (
	"errors"
	"fmt"
)

// ErrException indicates an exceptional condition occurred. For example, we cannot
// create an [*http.Request] because the URL is invalid. We use this wrapper error to
// distinguish between measurement errors and fundamental errors.
type ErrException struct {
	Err error
}

// Unwrap supports [errors.Unwrap].
func (exc *ErrException) Unwrap() error {
	return exc.Err
}

// Error implements error.
func (exc *ErrException) Error() string {
	return fmt.Sprintf("minilang: exception: %s", exc.Err.Error())
}

// IsErrException returns true when an error is an [ErrException].
func IsErrException(err error) bool {
	var exc *ErrException
	return errors.As(err, &exc)
}

// ErrSkip is a sentinel error indicating to a [Stage] that it should not run.
var ErrSkip = errors.New("minilang: skip this stage")

// IsErrSkip returns true when an error is an [ErrSkip].
func IsErrSkip(err error) bool {
	return errors.Is(err, ErrSkip)
}

// NewErrException creates a new exception with a formatted string as value.
func NewErrException(format string, v ...any) *ErrException {
	return &ErrException{fmt.Errorf(format, v...)}
}

// NewTypeErrException creates a new exception for the given types.
func NewTypeErrException[Expected any](got any) *ErrException {
	var expected Expected
	return NewErrException("type error: expected %T; got %T", expected, got)
}

// ErrInvalidNumberOfChildren indicates that the AST contains an invalid number of children.
var ErrInvalidNumberOfChildren = errors.New("dsl: invalid number of children")

// ErrUnwrap indicates that we cannot unwrap a [RunnableASTNode] to a [Stage] because
// there is a mismatch of the underlying [Stage] type.
var ErrUnwrap = errors.New("dsl: cannot unwrap RunnableASTNode to Stage")

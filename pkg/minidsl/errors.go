package minidsl

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
	return fmt.Sprintf("minidsl: exception: %s", exc.Err.Error())
}

// IsErrException returns true when an error is an [ErrException].
func IsErrException(err error) bool {
	var exc *ErrException
	return errors.As(err, &exc)
}

// ErrSkip is a sentinel error indicating to a stage that it should not run.
var ErrSkip = errors.New("minidsl: skip this stage")

// IsErrSkip returns true when an error is an [ErrSkip].
func IsErrSkip(err error) bool {
	return errors.Is(err, ErrSkip)
}

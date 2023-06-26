package dsl

import "fmt"

// Exception indicates that a [Function] generated an exception.
type Exception struct {
	Reason string
}

// NewException creates a new exception with the given reason string.
func NewException(format string, v ...any) *Exception {
	return &Exception{fmt.Sprintf(format, v...)}
}

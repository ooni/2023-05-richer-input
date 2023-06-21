package dslx

//
// Type management funcs
//

import "fmt"

// TypeString returns the string representation of a type.
func TypeString[T any]() string {
	var value T
	return fmt.Sprintf("%T", value)
}

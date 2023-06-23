package dsl

import "math"

// ExpectSingleScalarArgument expects the arguments to contain just a single value of type T.
func ExpectSingleScalarArgument[T any](arguments []any) (T, error) {
	if len(arguments) != 1 {
		return *new(T), NewErrCompile("expected single argument, got %T (%v)", arguments, arguments)
	}
	value, ok := arguments[0].(T)
	if !ok {
		return *new(T), NewErrCompile("cannot convert %T (%v) to %T", arguments[0], arguments[0], value)
	}
	return value, nil
}

// ExpectSingleUint16Argument is a specialization of ExpectSingleScalarArgument for uint16.
func ExpectSingleUint16Argument(arguments []any) (uint16, error) {
	// try handling the case where we set the uint16 explicitly first
	if value, err := ExpectSingleScalarArgument[uint16](arguments); err == nil {
		return value, nil
	}

	// fallback to the case where we did read the value from a JSON, which
	// will always represent numbers as float64 when not given a schema
	value, err := ExpectSingleScalarArgument[float64](arguments)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value >= math.MaxUint16 || math.Mod(value, 1) != 0 {
		return 0, NewErrCompile("cannot convert %T (%v) to uint16", value, value)
	}
	return uint16(math.Trunc(value)), nil
}

// ExpectListArguments expects the arguments to be a list of T.
func ExpectListArguments[T any](arguments []any) ([]T, error) {
	out := []T{}
	for _, argument := range arguments {
		value, ok := argument.(T)
		if !ok {
			return nil, NewErrCompile("cannot convert %T (%v) to %T", argument, argument, value)
		}
		out = append(out, value)
	}
	return out, nil
}

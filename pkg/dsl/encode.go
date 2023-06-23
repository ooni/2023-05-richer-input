package dsl

// EncodeNameProvider provides the name for EncodeXXX functions.
type EncodeNameProvider interface {
	Name() string
}

// EncodeFunctionScalar encodes a function with the name provided by the given
// name provider and whose argument must be a scalar with the given type.
func EncodeFunctionScalar[T any](provider EncodeNameProvider, value T) []any {
	return []any{provider.Name(), value}
}

// EncodeFunctionList is like EncodeFunctionScalar but for a list of arguments.
func EncodeFunctionList[T any](provider EncodeNameProvider, value []T) []any {
	out := []any{provider.Name()}
	for _, v := range value {
		out = append(out, v)
	}
	return out
}

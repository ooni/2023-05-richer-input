package dslx

// Void is an empty structure.
type Void struct{}

// IsVoid returns whether a type is [Void].
func IsVoid[T any]() bool {
	var value T
	switch (any)(value).(type) {
	case *Void:
		return true
	default:
		return false
	}
}

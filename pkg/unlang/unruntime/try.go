package unruntime

// Try checks whether the given value is an [*Exception] and unwraps the [*Exception] error
// in such a case. Otherwise, this [Try] returns nil to the caller.
func Try(value any) error {
	switch xvalue := value.(type) {
	case *Exception:
		return xvalue.Error

	default:
		return nil
	}
}

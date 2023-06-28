package rix

// Skip is a special type indicating to a [Func] that it should skip running. The [Func]
// must immediately return the [Skip] value to the caller when its input is a [Skip].
type Skip struct{}

package dslx

//
// Async helper functions
//

// asyncCollect collects all the elements returned by a channel under the
// assumption that the channel will be closed to signal EOF.
func asyncCollect[T any](c <-chan T) (v []T) {
	for t := range c {
		v = append(v, t)
	}
	return
}

// asyncStream returns a channel where it posts all the given elements and then
// closes the channel when it has finished posting the elements.
func asyncStream[T any](ts ...T) <-chan T {
	c := make(chan T)
	go func() {
		defer close(c) // as documented
		for _, t := range ts {
			c <- t
		}
	}()
	return c
}

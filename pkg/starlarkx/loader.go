package starlarkx

import (
	"errors"

	"go.starlark.net/starlark"
)

// Loader implements the `load` built-in function for starlark scripts. The zero value
// of this struct is invalid. Please, construct using the [NewLoader] function.
//
// The implementation of this struct derives from the official starlark-go implementation and
// therefore has the following license:
//
//	SPDX-License-Identifier: BSD-3-Clause
//
// The loader implementation from which this code derives is sequential.
type Loader struct {
	cache map[string]*loaderEntry
}

// NewLoader creates a new [*Loader] instance.
func NewLoader() *Loader {
	return &Loader{
		cache: map[string]*loaderEntry{},
	}
}

// loaderEntry is an entry in the [*Loader] cache.
type loaderEntry struct {
	globals starlark.StringDict
	err     error
}

// ErrCycle indicates there's a cycle in the load graph.
var ErrCycle = errors.New("starlarkx: cycle in load graph")

// Load implements the `load` built-in function.
func (lo *Loader) Load(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	entry, ok := lo.cache[module]
	if entry == nil {
		if ok {
			// request for package whose loading is in progress
			return nil, ErrCycle
		}

		// Add a placeholder to indicate "load in progress".
		lo.cache[module] = nil

		// Load it.
		thread := &starlark.Thread{
			Name:       "exec " + module,
			Print:      thread.Print,
			Load:       thread.Load,
			OnMaxSteps: thread.OnMaxSteps,
			Steps:      0,
		}
		predeclared := NewPredeclared()
		globals, err := starlark.ExecFile(thread, module, nil, predeclared)
		entry = &loaderEntry{globals, err}

		// Update the cache.
		lo.cache[module] = entry
	}
	return entry.globals, entry.err
}

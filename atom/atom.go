package atom

import (
	"sync"
)

// Atom represents an atomic value.
type Atom struct {
	v interface{}
	m sync.Mutex
}

// New creates a new Atom with the specified
// value.
func New(value interface{}) *Atom {
	return &Atom{v: value}
}

// Value gets the current value from the Atom.
func (a *Atom) Value() interface{} {
	a.m.Lock()
	defer a.m.Unlock()
	return a.v
}

// Set sets the specified value to the Atom.
func (a *Atom) Set(value interface{}) {
	a.m.Lock()
	defer a.m.Unlock()
	a.v = value
}

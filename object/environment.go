package object

import "sync"

type Environment struct {
	mu     sync.RWMutex
	values map[string]Object
	outer  *Environment
}

func NewEnvironment() *Environment {
	v := make(map[string]Object)
	return &Environment{values: v, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	obj, ok := e.values[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.values[name] = val
	return true
}

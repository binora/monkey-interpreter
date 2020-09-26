package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: map[string]Object{},
		outer: nil,
	}
}

func (e *Environment) Get(key string) (Object, bool) {
	val, ok := e.store[key]
	if !ok && e.outer != nil {
		return e.outer.Get(key)
	}
	return val, ok
}

func (e *Environment) Set(key string, value Object) Object {
	e.store[key] = value
	return value
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

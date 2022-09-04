package main

import "time"

type Callable interface {
	Call(Visitor, []interface{}) (interface{}, LoxError)

	Arity() int
}

func native() map[string]interface{} {
	return map[string]interface{}{
		"clock": Clock{},
	}
}

type Clock struct{}

func (c Clock) Arity() int {
	return 0
}

func (c Clock) Call(interpreter Visitor, arguments []interface{}) (interface{}, LoxError) {
	return time.Now().Unix(), nil
}

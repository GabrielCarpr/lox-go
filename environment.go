package main

import "fmt"

func NewGlobalEnvironment() *Environment {
	return &Environment{nil, native()}
}

func NewScopedEnvironment(from *Environment) *Environment {
	return &Environment{from, make(map[string]interface{})}
}

type Environment struct {
	enclosed *Environment
	values   map[string]interface{}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Assign(name Token, value interface{}) LoxError {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosed != nil {
		return e.enclosed.Assign(name, value)
	}

	return RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}

func (e Environment) Get(name Token) (interface{}, LoxError) {
	value, ok := e.values[name.lexeme]
	if ok {
		return value, nil
	}

	if e.enclosed != nil {
		return e.enclosed.Get(name)
	}

	return nil, RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}

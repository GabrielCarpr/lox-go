package main

import "fmt"

func NewEnvironment() *Environment {
	return &Environment{make(map[string]interface{})}
}

type Environment struct {
	values map[string]interface{}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Assign(name Token, value interface{}) LoxError {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}
	return RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
}

func (e Environment) Get(name Token) (interface{}, LoxError) {
	value, ok := e.values[name.lexeme]
	if !ok {
		return nil, RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
	}
	return value, nil
}

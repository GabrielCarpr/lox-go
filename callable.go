package main

import (
	"fmt"
	"time"
)

type Callable interface {
	Call(*Interpreter, []interface{}) (interface{}, LoxError)

	Arity() int

	String() string
}

var _ Callable = &LoxFunction{}

type LoxFunction struct {
	declaration Function
}

func (f *LoxFunction) Call(interpreter *Interpreter, args []interface{}) (interface{}, LoxError) {
	environment := NewScopedEnvironment(interpreter.globals)

	for i, param := range f.declaration.Params {
		environment.Define(param.lexeme, args[i])
	}

	err := interpreter.executeBlock(f.declaration.Body.Statements, environment)
	return nil, err
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.lexeme)
}

// Native functions

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

func (c Clock) String() string {
	return "<native fn>"
}

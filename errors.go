package main

type LoxError interface {
	Type() string
	Token() Token
	error
}

type RuntimeError struct {
	SubjectToken Token
	Message      string
}

func (e RuntimeError) Type() string {
	return "RuntimeError"
}

func (e RuntimeError) Token() Token {
	return e.SubjectToken
}

func (e RuntimeError) Error() string {
	return e.Message
}

type DivideZeroError struct {
	RuntimeError
}

func (e DivideZeroError) Type() string {
	return "DivideZeroError"
}

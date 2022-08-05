package main

import "fmt"

type Lexeme int

const (
	LEFT_PAREN Lexeme = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL
	IDENTIFIER
	STRING
	NUMBER
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

type Token struct {
	tokenType Lexeme
	lexeme    string
	literal   string
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %d %s", t.tokenType, t.lexeme, t.literal)
}

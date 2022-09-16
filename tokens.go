package main

import "fmt"

type Lexeme string

const (
	LEFT_PAREN    Lexeme = "LEFT_PAREN"
	RIGHT_PAREN          = "RIGHT_PAREN"
	LEFT_BRACE           = "LEFT_BRACE"
	RIGHT_BRACE          = "RIGHT_BRACE"
	COMMA                = "COMMA"
	DOT                  = "DOT"
	MINUS                = "MINUS"
	PLUS                 = "PLUS"
	SEMICOLON            = "SEMICOLON"
	SLASH                = "SLASH"
	STAR                 = "STAR"
	BANG                 = "BANG"
	BANG_EQUAL           = "BANG_EQUAL"
	EQUAL                = "EQUAL"
	EQUAL_EQUAL          = "EQUAL_EQUAL"
	GREATER              = "GREATER"
	GREATER_EQUAL        = "GREATER_EQUAL"
	LESS                 = "LESS"
	LESS_EQUAL           = "LESS_EQUAL"
	IDENTIFIER           = "IDENTIFIER"
	STRING               = "STRING"
	NUMBER               = "NUMBER"
	AND                  = "AND"
	CLASS                = "CLASS"
	ELSE                 = "ELSE"
	FALSE                = "FALSE"
	FUN                  = "FUN"
	FOR                  = "FOR"
	IF                   = "IF"
	NIL                  = "NIL"
	OR                   = "OR"
	PRINT                = "PRINT"
	RETURN               = "RETURN"
	SUPER                = "SUPER"
	THIS                 = "THIS"
	TRUE                 = "TRUE"
	VAR                  = "VAR"
	WHILE                = "WHILE"
	PERCENT              = "PERCENT"
	EOF                  = "EOF"
)

type Token struct {
	tokenType Lexeme
	lexeme    string
	literal   interface{}
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %s Lexeme: %s Literal: %v", t.tokenType, t.lexeme, t.literal)
}

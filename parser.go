package main

import (
	"errors"
	"fmt"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) Load(tokens []Token) {
	p.tokens = tokens
	p.current = 0
}

func (p *Parser) Parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{false}
	}

	if p.match(TRUE) {
		return Literal{true}
	}

	if p.match(NIL) {
		return Literal{nil}
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}

	panic(p.error(p.previous(), "Unexpected token"))
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.advance()
	}
}

// Parsing infrastructure

func (p *Parser) consume(tokenType Lexeme, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(p.error(p.previous(), message))
}

func (p *Parser) error(token Token, message string) error {
	var err error
	if token.tokenType == EOF {
		err = errors.New(fmt.Sprintf("%s at end", message))
	} else {
		err = errors.New(fmt.Sprintf("%s at %s", message, token.lexeme))
	}

	report(err, token.line)
	return err
}

func (p *Parser) match(types ...Lexeme) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType Lexeme) bool {
	return p.peek().tokenType == tokenType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
		return p.previous()
	}
	return p.peek()
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

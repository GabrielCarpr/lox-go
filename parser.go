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

func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return []Stmt{}, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

// Grammar

func (p *Parser) declaration() (Stmt, error) {
	var result Stmt
	var err error
	if p.match(VAR) {
		result, err = p.varDeclaration()
	} else {
		result, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil, err
	}
	return result, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expected variable name")
	if err != nil {
		return nil, err
	}

	var init Expr
	if p.match(EQUAL) {
		init, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "';' expected after variable declaration")
	if err != nil {
		return nil, err
	}
	return Var{name, &init}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return p.block()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "';' expected after print value")
	return Print{value}, err
}

func (p *Parser) block() (Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, dec)
	}

	_, err := p.consume(RIGHT_BRACE, "'}' expected after block")
	return Block{statements}, err
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "'(' expected after 'if'")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "')' expected after 'if' condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return If{condition, thenBranch, elseBranch}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "'(' expected after while statement")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "')' expected after 'while' condition")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	return While{condition, body}, err
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "';' expected after value")
	return Expression{value}, err
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if assign, ok := expr.(Variable); ok {
			name := assign.Name
			return Assign{name, value}, nil
		}
		return nil, p.error(equals, "Invalid assignment target")
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		return Logical{expr, operator, right}, err
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		return Logical{expr, operator, right}, err
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{operator, right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return Literal{false}, nil
	}

	if p.match(TRUE) {
		return Literal{true}, nil
	}

	if p.match(NIL) {
		return Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression")
		return Grouping{expr}, err
	}

	if p.match(IDENTIFIER) {
		return Variable{p.previous()}, nil
	}

	return nil, p.error(p.previous(), "Unexpected token")
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

func (p *Parser) consume(until Lexeme, message string) (Token, error) {
	if p.check(until) {
		return p.advance(), nil
	}

	return Token{}, p.error(p.previous(), message)
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

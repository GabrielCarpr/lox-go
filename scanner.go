package main

import (
	"fmt"
)

func GetScanner(source string) *Scanner {
	return &Scanner{source, 0, 1, 0, false}
}

type Scanner struct {
	source        string
	start         int
	line          int
	current       int
	containsError bool
}

func (s *Scanner) Scan() ([]Token, error) {
	tokens := make([]Token, 0)

	for !s.isAtEnd() {
		token, err := s.tokenize()
		if err != nil {
			report(err, s.line)
			s.containsError = true
		}
		tokens = append(tokens, token)
	}

	tokens = append(tokens, Token{
		tokenType: EOF,
		lexeme:    "",
		literal:   "",
		line:      s.line,
	})

	return tokens, nil
}

func (s *Scanner) tokenize() (Token, error) {
	lex, err := s.scan()
	if err != nil {
		return Token{}, err
	}

	text := s.source[s.start:s.current]

	return Token{
		tokenType: lex,
		lexeme:    text,
		literal:   "",
		line:      s.line,
	}, nil
}

func (s *Scanner) scan() (Lexeme, error) {
	char := s.advance()

	var lex Lexeme
	switch char {
	case '(':
		lex = LEFT_PAREN
	case ')':
		lex = RIGHT_PAREN
	case '{':
		lex = LEFT_PAREN
	case '}':
		lex = RIGHT_PAREN
	case ',':
		lex = COMMA
	case '.':
		lex = DOT
	case '-':
		lex = MINUS
	case '+':
		lex = PLUS
	case ';':
		lex = SEMICOLON
	case '*':
		lex = STAR
	default:
		return 0, fmt.Errorf("unexpected character: %s", string(char))
	}

	return lex, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current])
}

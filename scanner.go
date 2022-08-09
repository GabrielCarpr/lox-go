package main

import (
	"errors"
	"fmt"
	"strconv"
)

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:        source,
		start:         0,
		line:          1,
		current:       0,
		containsError: false,
		tokens:        make([]Token, 0)}
}

type Scanner struct {
	source        string
	start         int
	line          int
	current       int
	containsError bool
	tokens        []Token
}

func (s *Scanner) Scan() {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scan()
		if err != nil {
			report(err, s.line)
			s.containsError = true
		}
	}

	s.tokenize(EOF, nil)
}

func (s *Scanner) tokenize(lex Lexeme, literal interface{}) error {
	text := s.source[s.start:s.current]

	s.tokens = append(s.tokens, Token{
		tokenType: lex,
		lexeme:    text,
		literal:   literal,
		line:      s.line,
	})

	return nil
}

func (s *Scanner) scan() error {
	char := s.advance()

	switch char {
	case "(":
		s.tokenize(LEFT_PAREN, nil)
		break
	case ")":
		s.tokenize(RIGHT_PAREN, nil)
		break
	case "{":
		s.tokenize(LEFT_PAREN, nil)
		break
	case "}":
		s.tokenize(RIGHT_PAREN, nil)
		break
	case ",":
		s.tokenize(COMMA, nil)
		break
	case ".":
		s.tokenize(DOT, nil)
		break
	case "-":
		s.tokenize(MINUS, nil)
		break
	case "+":
		s.tokenize(PLUS, nil)
		break
	case ";":
		s.tokenize(SEMICOLON, nil)
		break
	case "*":
		s.tokenize(STAR, nil)
		break
	case "!":
		if s.matchNext("=") {
			s.tokenize(BANG_EQUAL, nil)
		} else {
			s.tokenize(BANG, nil)
		}
		break
	case "=":
		if s.matchNext("=") {
			s.tokenize(EQUAL_EQUAL, nil)
		} else {
			s.tokenize(EQUAL, nil)
		}
	case "<":
		if s.matchNext("=") {
			s.tokenize(LESS_EQUAL, nil)
		} else {
			s.tokenize(LESS, nil)
		}
		break
	case ">":
		if s.matchNext("=") {
			s.tokenize(GREATER_EQUAL, nil)
		} else {
			s.tokenize(GREATER, nil)
		}
		break
	case "/":
		if s.matchNext("/") {
			for s.peek(0) != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.tokenize(SLASH, nil)
		}
		break
	case " ":
	case "\r":
	case "\t":
	case "":
		break

	case "\n":
		s.line++
		break

	case "\"":
		s.scanString()
		break

	default:
		if isDigit(char) {
			s.scanNumber()
			break
		}

		return fmt.Errorf("unexpected character: %s", string(char))
	}

	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() string {
	s.current++
	return string(s.source[s.current-1])
}

func (s *Scanner) matchNext(expected string) bool {
	if s.peek(0) != expected {
		return false
	}

	s.advance()
	return true
}

func (s *Scanner) peek(next int) string {
	if s.isAtEnd() {
		return "0"
	}
	return string(s.source[s.current+next])
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek(0)) {
		s.advance()
	}

	if s.peek(0) == "." && isDigit(s.peek(1)) {
		s.advance()
	}

	for isDigit(s.peek(0)) {
		s.advance()
	}

	literal, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		report(errors.New("could not parse number"), s.line)
		s.containsError = true
		return
	}

	s.tokenize(NUMBER, literal)
}

func (s *Scanner) scanString() {
	for s.peek(0) != "\"" && !s.isAtEnd() {
		if s.peek(0) == "\n" {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		report(errors.New("unterminated string"), s.line)
		s.containsError = true
		return
	}
	s.advance()
	s.tokenize(STRING, s.source[s.start+1:s.current-1])
}

func isDigit(char string) bool {
	_, err := strconv.ParseInt(char, 10, 64)
	return err == nil
}

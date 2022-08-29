package main

func NewLox() *Lox {
	return &Lox{
		NewParser(),
		NewInterpreter(),
	}
}

type Lox struct {
	parser      *Parser
	interpreter *Interpreter
}

func (l *Lox) Run(source string) error {
	scanner := NewScanner(source)
	scanner.Scan()
	tokens := scanner.tokens

	l.parser.Load(tokens)
	ast, err := l.parser.Parse()
	if err != nil {
		return err
	}

	return l.interpreter.Interpret(ast)
}

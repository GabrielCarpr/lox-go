package main

func NewLox() *Lox {
	interpreter := NewInterpreter()
	return &Lox{
		NewParser(),
		interpreter,
	}
}

type Lox struct {
	parser      *Parser
	interpreter *Interpreter
}

func (l *Lox) Run(source string) error {
	scanner := NewScanner(source)
	resolver := NewResolver(l.interpreter)
	scanner.Scan()
	tokens := scanner.tokens

	l.parser.Load(tokens)
	ast, err := l.parser.Parse()
	if err != nil {
		return err
	}

	resolver.Resolve(ast)
	return l.interpreter.Interpret(ast)
}

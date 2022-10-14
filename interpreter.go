package main

import (
	"fmt"
	"regexp"
)

func NewInterpreter() *Interpreter {
	globals := NewGlobalEnvironment()

	return &Interpreter{globals, globals, make(map[Expr]int64)}
}

var _ Visitor = (&Interpreter{})

type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[Expr]int64
}

func (i *Interpreter) Interpret(statements []Stmt) error {
	for _, stmt := range statements {
		err := i.execute(stmt)
		if err != nil {
			report(err, err.Token().line)
			return err
		}
	}
	return nil
}

// Visitor methods

func (i *Interpreter) VisitPrintStmt(stmt Print) LoxError {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return err
	}

	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt Expression) LoxError {
	_, err := i.evaluate(stmt.Expression)
	return err
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) (interface{}, LoxError) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) (interface{}, LoxError) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) (interface{}, LoxError) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.tokenType {
	case MINUS:
		num, err := checkNumber(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -num, nil
	case BANG:
		return !i.isTruthy(right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr Call) (interface{}, LoxError) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	arguments := make([]interface{}, len(expr.Arguments))
	for j, arg := range expr.Arguments {
		argument, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments[j] = argument
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, RuntimeError{expr.Paren, fmt.Sprintf("Can only call functions and classes")}
	}
	if len(arguments) != function.Arity() {
		return nil, RuntimeError{expr.Paren, fmt.Sprintf("Expected %d arguments, got %d", function.Arity(), len(arguments))}
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) VisitVarStmt(stmt Var) LoxError {
	if stmt.Initialiser != nil {
		value, err := i.evaluate(*stmt.Initialiser)
		if err != nil {
			return err
		}
		i.environment.Define(stmt.Name.lexeme, value)
	} else {
		i.environment.Define(stmt.Name.lexeme, nil)
	}
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt Block) LoxError {
	return i.executeBlock(stmt.Statements, NewScopedEnvironment(i.environment))
}

func (i *Interpreter) VisitIfStmt(stmt If) LoxError {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	if i.isTruthy(condition) {
		if err := i.execute(stmt.Then); err != nil {
			return err
		}
	} else if stmt.Else != nil {
		if err := i.execute(stmt.Else); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitFunctionStmt(fn Function) LoxError {
	function := &LoxFunction{fn, i.environment}
	i.environment.Define(fn.Name.lexeme, function)
	return nil
}

func (i *Interpreter) VisitReturnStmt(ret Return) LoxError {
	var value interface{}
	var err LoxError
	if ret.Value != nil {
		value, err = i.evaluate(*ret.Value)
		if err != nil {
			return err
		}
	}

	return ReturnError{value, ret.Keyword}
}

func (i *Interpreter) VisitVariableExpr(expr Variable) (interface{}, LoxError) {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) (interface{}, LoxError) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.tokenType {
	case MINUS:
		l, r, err := checkNumbers(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l - r, nil
	case SLASH:
		l, r, err := checkNumbers(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		if r == 0 {
			return nil, DivideZeroError{RuntimeError{expr.Operator, "Cannot divide by zero"}}
		}
		return l / r, nil
	case STAR:
		l, r, err := checkNumbers(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return l * r, nil
	case PLUS:
		leftNum, isLeftNum := left.(float64)
		rightNum, isRightNum := right.(float64)
		leftString, isLeftString := left.(string)
		rightString, isRightString := right.(string)

		if isLeftNum && isRightNum {
			return leftNum + rightNum, nil
		}

		if isLeftString && isRightString {
			return leftString + rightString, nil
		}

		return nil, &RuntimeError{expr.Operator, "Operands must be two strings or two numbers"}
	case GREATER:
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case LESS:
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	case PERCENT:
		leftNum, isLeftNum := left.(float64)
		rightNum, isRightNum := right.(float64)
		if !isLeftNum || !isRightNum {
			return nil, RuntimeError{expr.Operator, "Operands must be two numbers"}
		}
		return int64(leftNum) % int64(rightNum), nil
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr Assign) (interface{}, LoxError) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[expr]
	if ok {
		i.environment.AtDepth(distance).Assign(expr.Name, value)
	}

	return value, err
}

func (i *Interpreter) VisitLogicalExpr(expr Logical) (interface{}, LoxError) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.tokenType == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitWhileStmt(stmt While) LoxError {
	for true {
		condition, err := i.evaluate(stmt.Condition)
		if err != nil {
			return err
		}
		if truthy, ok := condition.(bool); !truthy || !ok {
			break
		}

		err = i.execute(stmt.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

// Private methods

func (i *Interpreter) execute(stmt Stmt) LoxError {
	return stmt.Accept(i)
}

func (i *Interpreter) resolve(expr Expr, depth int64) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(stmts []Stmt, env *Environment) LoxError {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = env

	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, LoxError) {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}

	switch t := obj.(type) {
	case bool:
		return t
	}

	return false
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return a == b
}

func (i *Interpreter) lookupVariable(name Token, expr Expr) (interface{}, LoxError) {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.AtDepth(distance).Get(name)
	} else {
		return i.globals.Get(name)
	}
}

// Utility methods

func stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	if str, ok := obj.(string); ok {
		return str
	}

	if num, ok := obj.(float64); ok {
		text := fmt.Sprintf("%f", num)
		reggie := regexp.MustCompile("^(.*)\\.0+$")
		if reggie.MatchString(text) {
			text = reggie.FindStringSubmatch(text)[1]
		}
		reggie2 := regexp.MustCompile("^(.*\\.[1-9]+)0+")
		if reggie2.MatchString(text) {
			text = reggie2.FindStringSubmatch(text)[1]
		}
		return text
	}

	if str, ok := obj.(fmt.Stringer); ok {
		return str.String()
	}

	return fmt.Sprintf("%v+", obj)
}

func checkNumber(operator Token, value interface{}) (float64, *RuntimeError) {
	if number, ok := value.(float64); ok {
		return number, nil
	}
	return 0, &RuntimeError{operator, "Operand must be a number"}
}

func checkNumbers(operator Token, a interface{}, b interface{}) (aNum float64, bNum float64, err *RuntimeError) {
	aNum, aOk := a.(float64)
	bNum, bOk := b.(float64)

	if !aOk || !bOk {
		err = &RuntimeError{operator, "Must be a number"}
	}

	return aNum, bNum, nil
}

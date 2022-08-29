package main

import (
	"fmt"
	"regexp"
)

func NewInterpreter() *Interpreter {
	return &Interpreter{NewEnvironment()}
}

type Interpreter struct {
	environment *Environment
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

func (i *Interpreter) VisitVariableExpr(expr Variable) (interface{}, LoxError) {
	return i.environment.Get(expr.Name)
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
	}

	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr Assign) (interface{}, LoxError) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	err = i.environment.Assign(expr.Name, value)
	return value, err
}

// Private methods

func (i *Interpreter) execute(stmt Stmt) LoxError {
	return stmt.Accept(i)
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

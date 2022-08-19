package main

import "fmt"

type ASTPrinter struct {
}

func (p *ASTPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *ASTPrinter) VisitBinaryExpr(expr Binary) interface{} {
	return p.parenthesize(expr.Operator.lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) VisitGroupingExpr(expr Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *ASTPrinter) VisitUnaryExpr(expr Unary) interface{} {
	return p.parenthesize(expr.Operator.lexeme, expr.Right)
}

func (p *ASTPrinter) VisitLiteralExpr(expr Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
	result := ""
	result += "("
	result += name

	for _, expr := range exprs {
		result += " "
		result += fmt.Sprintf("%s", expr.Accept(p))
	}

	result += ")"

	return result
}

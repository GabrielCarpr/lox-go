package main

type Visitor interface {
	VisitBinaryExpr(expr Binary) interface{}
	VisitUnaryExpr(expr Unary) interface{}
	VisitGroupingExpr(expr Grouping) interface{}
	VisitLiteralExpr(expr Literal) interface{}
}

type Expr interface {
	Accept(Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(b)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(u)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(l)
}

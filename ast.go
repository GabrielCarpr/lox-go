package main

type ExprVisitor[R any] interface {
	VisitBinary(Binary) R
}

type Node struct {
}

func (e Node) Accept(visitor ExprVisitor) {

}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value interface{}
}

type Unary struct {
	Operator Token
	Right    Expr
}

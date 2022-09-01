package main

// Expressions
type ExprVisitor interface {
	VisitBinaryExpr(expr Binary) (interface{}, LoxError)
	VisitUnaryExpr(expr Unary) (interface{}, LoxError)
	VisitGroupingExpr(expr Grouping) (interface{}, LoxError)
	VisitLiteralExpr(expr Literal) (interface{}, LoxError)
	VisitVariableExpr(expr Variable) (interface{}, LoxError)
	VisitAssignExpr(expr Assign) (interface{}, LoxError)
}

type Expr interface {
	Accept(ExprVisitor) (interface{}, LoxError)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b Binary) Accept(visitor ExprVisitor) (interface{}, LoxError) {
	return visitor.VisitBinaryExpr(b)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (u Unary) Accept(visitor ExprVisitor) (interface{}, LoxError) {
	return visitor.VisitUnaryExpr(u)
}

type Grouping struct {
	Expression Expr
}

func (g Grouping) Accept(visitor ExprVisitor) (interface{}, LoxError) {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor ExprVisitor) (interface{}, LoxError) {
	return visitor.VisitLiteralExpr(l)
}

type Variable struct {
	Name Token
}

func (v Variable) Accept(visitor ExprVisitor) (interface{}, LoxError) {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Name  Token
	Value Expr
}

func (a Assign) Accept(visitor ExprVisitor) (interface{}, LoxError) {
	return visitor.VisitAssignExpr(a)
}

// Statements

type StmtVisitor interface {
	VisitExpressionStmt(Expression) LoxError
	VisitPrintStmt(Print) LoxError
	VisitVarStmt(Var) LoxError
	VisitBlockStmt(Block) LoxError
}

type Stmt interface {
	Accept(StmtVisitor) LoxError
}

type Expression struct {
	Expression Expr
}

func (e Expression) Accept(visitor StmtVisitor) LoxError {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	Expression Expr
}

func (p Print) Accept(visitor StmtVisitor) LoxError {
	return visitor.VisitPrintStmt(p)
}

type Var struct {
	Name        Token
	Initialiser *Expr
}

func (v Var) Accept(visitor StmtVisitor) LoxError {
	return visitor.VisitVarStmt(v)
}

type Block struct {
	Statements []Stmt
}

func (b Block) Accept(visitor StmtVisitor) LoxError {
	return visitor.VisitBlockStmt(b)
}

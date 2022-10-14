package main

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter, make([]map[string]bool, 0)}
}

var _ Visitor = &Resolver{}

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func (r *Resolver) Resolve(ast []Stmt) {
	r.resolveStmt(ast...)
}

func (r *Resolver) resolveStmt(stmts ...Stmt) {
	for _, stmt := range stmts {
		stmt.Accept(r)
	}
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[0 : len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.lexeme] = true
}

// Visit methods

func (r *Resolver) VisitBlockStmt(block Block) LoxError {
	r.beginScope()
	r.resolveStmt(block.Statements...)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVarStmt(v Var) LoxError {
	r.declare(v.Name)
	if v.Initialiser != nil {
		r.resolveExpr(*v.Initialiser)
	}
	r.define(v.Name)

	return nil
}

func (r *Resolver) VisitVariableExpr(v Variable) (interface{}, LoxError) {
	hasScopes := len(r.scopes) != 0
	if hasScopes {
		if ready, ok := r.scopes[len(r.scopes)-1][v.Name.lexeme]; ok && !ready {
			return nil, CompileError{v.Name, "Cannot read local variable in it's own initializer"}
		}
	}

	r.resolveLocal(v, v.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr Assign) (interface{}, LoxError) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(fun Function) LoxError {
	r.declare(fun.Name)
	r.define(fun.Name)

	r.resolveFunction(fun)
	return nil
}

func (r *Resolver) VisitExpressionStmt(expr Expression) LoxError {
	r.resolveStmt(expr)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt If) LoxError {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Then)
	if stmt.Else != nil {
		r.resolveStmt(stmt.Else)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(print Print) LoxError {
	r.resolveExpr(print.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(ret Return) LoxError {
	if ret.Value != nil {
		r.resolveExpr(*ret.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt While) LoxError {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitBinaryExpr(bin Binary) (interface{}, LoxError) {
	r.resolveExpr(bin.Left)
	r.resolveExpr(bin.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(call Call) (interface{}, LoxError) {
	r.resolveExpr(call.Callee)
	for _, arg := range call.Arguments {
		r.resolveExpr(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(group Grouping) (interface{}, LoxError) {
	r.resolveExpr(group.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(lit Literal) (interface{}, LoxError) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(log Logical) (interface{}, LoxError) {
	r.resolveExpr(log.Left)
	r.resolveExpr(log.Right)

	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(unary Unary) (interface{}, LoxError) {
	r.resolveExpr((unary.Right))
	return nil, nil
}

// Helpers

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i, scope := range r.scopes {
		if _, ok := scope[name.lexeme]; ok {
			r.interpreter.resolve(expr, int64(len(r.scopes)-1-i))
		}
	}
}

func (r *Resolver) resolveFunction(fun Function) {
	r.beginScope()
	for _, param := range fun.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmt(fun.Body)
	r.endScope()
}

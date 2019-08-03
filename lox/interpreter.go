package lox

import (
	"fmt"
	"reflect"
)

// Interpreter ...
type Interpreter struct {
	Env       *Environment
	GlobalEnv *Environment
}

// NewInterpreter ...
func NewInterpreter() *Interpreter {
	ni := new(Interpreter)
	ni.Env = NewEnvironment(nil)
	ni.GlobalEnv = ni.Env
	ni.GlobalEnv.Define("clock", &Clock{})
	return ni
}

// VisitLiteralExpr ...
func (i Interpreter) VisitLiteralExpr(expr *ExprLiteral) interface{} {
	return expr.Value
}

// VisitGroupingExpr ...
func (i Interpreter) VisitGroupingExpr(expr *ExprGrouping) interface{} {
	return i.evaluate(expr.Expr)
}

// VisitVarExpr ...
func (i Interpreter) VisitVarExpr(expr *ExprVar) interface{} {
	return i.Env.Get(expr.Name.Lexeme)
}

// VisitLogicalExpr ...
func (i Interpreter) VisitLogicalExpr(expr *ExprLogical) interface{} {
	left := i.evaluate(expr.Left)
	if expr.Operator.Type == TokenTypeOr {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.Right)
}

// Interpret ...
func (i Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i Interpreter) execute(stmt Stmt) interface{} {
	return stmt.Accept(i)
}

func (i Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

// VisitAssignExpr ...
func (i Interpreter) VisitAssignExpr(expr *ExprAssign) interface{} {
	value := i.evaluate(expr.Value)
	i.Env.Assign(expr.Name.Lexeme, value)
	return value
}

// VisitUnaryExpr ...
func (i Interpreter) VisitUnaryExpr(expr *ExprUnary) interface{} {
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case TokenTypeBang:
		return !i.isTruthy(right.(float64))
	case TokenTypeMinus:
		return -right.(float64)
	}
	return nil
}

func (i Interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	switch obj.(type) {
	case string:
		return !(obj.(string) == "")
	case int, float64:
		return !(obj.(float64) == float64(0))
	case bool:
		return bool(obj.(bool))
	default:
		return true
	}
}

// VisitCallExpr ...
func (i Interpreter) VisitCallExpr(expr *ExprCall) interface{} {
	callee := i.evaluate(expr.Callee)
	arguments := make([]interface{}, 0)
	for _, arg := range expr.Arguments {
		arguments = append(arguments, i.evaluate(*arg))
	}
	f, ok := callee.(Callable)
	if ok {
		if len(arguments) != f.Arity() {
			fmt.Printf("Expected %d arguments but got %d\n", f.Arity(), len(arguments))
			return nil
		}
		return f.Call(&i, arguments)
	}
	fmt.Printf("Can only call functions and classes, not %v.\n", reflect.TypeOf(callee))
	return nil
}

// VisitBinaryExpr ...
func (i Interpreter) VisitBinaryExpr(expr *ExprBinary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case TokenTypeGreater:
		return left.(float64) > right.(float64)
	case TokenTypeGreaterEqual:
		return left.(float64) >= right.(float64)
	case TokenTypeLess:
		return left.(float64) < right.(float64)
	case TokenTypeLessEqual:
		return left.(float64) <= right.(float64)
	case TokenTypeBangEqual:
		return !(left.(float64) == right.(float64))
	case TokenTypeEqualEqual:
		return left.(float64) == right.(float64)
	case TokenTypeMinus:
		return left.(float64) - right.(float64)
	case TokenTypePlus:
		typeLeft := reflect.TypeOf(left).String()
		typeRight := reflect.TypeOf(right).String()
		if (typeLeft == "float64" || typeLeft == "float32" || typeLeft == "int") && (typeRight == "float64" || typeRight == "float32" || typeRight == "int") {
			return left.(float64) + right.(float64)
		} else if typeLeft == "string" && typeRight == "string" {
			return left.(string) + right.(string)
		}
	case TokenTypeSlash:
		return left.(float64) / right.(float64)
	case TokenTypeStar:
		return left.(float64) * right.(float64)
	}
	return nil
}

// VisitExpressionStmt ...
func (i Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) interface{} {
	return i.evaluate(stmt.Expression)
}

// VisitPrintStmt ...
func (i Interpreter) VisitPrintStmt(stmt *PrintStmt) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Printf("%v\n", value)
	return nil
}

// VisitVarStmt ...
func (i Interpreter) VisitVarStmt(stmt *VarStmt) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.Env.Define(stmt.Name.Lexeme, value)
	return nil
}

// VisitBlockStmt ...
func (i Interpreter) VisitBlockStmt(stmt *BlockStmt) interface{} {
	i.ExecuteBlock(stmt.Statements, NewEnvironment(i.Env))
	return nil
}

// ExecuteBlock ...
func (i Interpreter) ExecuteBlock(statements []Stmt, env *Environment) {
	previous := i.Env
	i.Env = env
	for _, statement := range statements {
		i.execute(statement)
	}
	i.Env = previous
}

// VisitIfStmt ...
func (i Interpreter) VisitIfStmt(stmt *IfStmt) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

// VisitWhileStmt ...
func (i Interpreter) VisitWhileStmt(stmt *WhileStmt) interface{} {
	for {
		if !i.isTruthy(i.evaluate(stmt.Condition)) {
			break
		}
		i.execute(stmt.Body)
	}
	return nil
}

// VisitFunctionStmt ...
func (i Interpreter) VisitFunctionStmt(stmt *FunctionStmt) interface{} {
	f := NewFunction(*stmt)
	i.Env.Define(stmt.Name.Lexeme, f)
	return nil
}

// VisitReturnStmt ...
func (i Interpreter) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	return &ReturnValue{ExprLiteral{Value: value}}
}

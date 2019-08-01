package lox

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	Env       *Environment
	GlobalEnv *Environment
}

func NewInterpreter() *Interpreter {
	ni := new(Interpreter)
	ni.Env = NewEnvironment(nil)
	ni.GlobalEnv = ni.Env
	ni.GlobalEnv.Define("clock", &Clock{})
	return ni
}

func (i Interpreter) VisitLiteralExpr(expr *ExprLiteral) interface{} {
	return expr.Value
}

func (i Interpreter) VisitGroupingExpr(expr *ExprGrouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i Interpreter) VisitVarExpr(expr *ExprVar) interface{} {
	return i.Env.Get(expr.Name.Lexeme)
}

func (i Interpreter) VisitLogicalExpr(expr *ExprLogical) interface{} {
	left := i.evaluate(expr.Left)
	if expr.Operator.Type == OR {
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

func (i Interpreter) VisitAssignExpr(expr *ExprAssign) interface{} {
	value := i.evaluate(expr.Value)
	i.Env.Assign(expr.Name.Lexeme, value)
	return value
}

func (i Interpreter) VisitUnaryExpr(expr *ExprUnary) interface{} {
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case BANG:
		return !i.isTruthy(right.(float64))
	case MINUS:
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
	} else {
		fmt.Printf("Can only call functions and classes, not %v.\n", reflect.TypeOf(callee))
		return nil
	}
}

func (i Interpreter) VisitBinaryExpr(expr *ExprBinary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !(left.(float64) == right.(float64))
	case EQUAL_EQUAL:
		return left.(float64) == right.(float64)
	case MINUS:
		return left.(float64) - right.(float64)
	case PLUS:
		typeLeft := reflect.TypeOf(left).String()
		typeRight := reflect.TypeOf(right).String()
		if (typeLeft == "float64" || typeLeft == "float32" || typeLeft == "int") && (typeRight == "float64" || typeRight == "float32" || typeRight == "int") {
			return left.(float64) + right.(float64)
		} else if typeLeft == "string" && typeRight == "string" {
			return left.(string) + right.(string)
		}
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	}
	return nil
}

func (i Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) interface{} {
	return i.evaluate(stmt.Expression)
}

func (i Interpreter) VisitPrintStmt(stmt *PrintStmt) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Printf("%v\n", value)
	return nil
}

func (i Interpreter) VisitVarStmt(stmt *VarStmt) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.Env.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i Interpreter) VisitBlockStmt(stmt *BlockStmt) interface{} {
	i.ExecuteBlock(stmt.Statements, NewEnvironment(i.Env))
	return nil
}

func (i Interpreter) ExecuteBlock(statements []Stmt, env *Environment) {
	previous := i.Env
	i.Env = env
	for _, statement := range statements {
		i.execute(statement)
	}
	i.Env = previous
}

func (i Interpreter) VisitIfStmt(stmt *IfStmt) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i Interpreter) VisitWhileStmt(stmt *WhileStmt) interface{} {
	for {
		if !i.isTruthy(i.evaluate(stmt.Condition)) {
			break
		}
		i.execute(stmt.Body)
	}
	return nil
}

func (i Interpreter) VisitFunctionStmt(stmt *FunctionStmt) interface{} {
	f := NewFunction(*stmt)
	i.Env.Define(stmt.Name.Lexeme, f)
	return nil
}

func (i Interpreter) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	return &ReturnValue{Value: value}
}

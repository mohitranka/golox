package lox

// Function ...
type Function struct {
	Declaration FunctionStmt
}

// NewFunction ...
func NewFunction(declaration FunctionStmt) *Function {
	return &Function{Declaration: declaration}
}

// Call ...
func (f Function) Call(i *Interpreter, args []interface{}) interface{} {
	env := NewEnvironment(i.GlobalEnv)
	for idx, param := range f.Declaration.Params {
		env.Define(param.Lexeme, args[idx])
	}
	i.ExecuteBlock(f.Declaration.Body, env)
	return nil
}

// Arity ...
func (f Function) Arity() int {
	return len(f.Declaration.Params)
}

package lox

type Function struct {
	Declaration FunctionStmt
}

func NewFunction(declaration FunctionStmt) *Function {
	return &Function{Declaration: declaration}
}

func (f Function) Call(i *Interpreter, args []interface{}) interface{} {
	env := NewEnvironment(i.GlobalEnv)
	for idx, param := range f.Declaration.Params {
		env.Define(param.Lexeme, args[idx])
	}
	i.ExecuteBlock(f.Declaration.Body, env)
	return nil
}

func (f Function) Arity() int {
	return len(f.Declaration.Params)
}

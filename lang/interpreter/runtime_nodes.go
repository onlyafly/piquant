package interpreter

import (
	"fmt"
	"piquant/lang/ast"
	"piquant/lang/token"
)

////////// EnvNode

type EnvNode struct {
	Env Env
}

func NewEnvNode(e Env) *EnvNode {
	return &EnvNode{Env: e}
}

func (en *EnvNode) Name() string {
	return en.Env.Name()
}

func (en *EnvNode) String() string {
	return "#environment<" + en.Env.Name() + ">"
}

func (en *EnvNode) isExpr() bool         { return true }
func (en *EnvNode) TypeName() string     { return "environment" }
func (en *EnvNode) Loc() *token.Location { return nil }
func (en *EnvNode) Equals(n ast.Node) bool {
	panicEvalError(n, "Cannot compare the values of environments: "+
		en.String()+" and "+n.String())
	return false
}

////////// Primitive

type primitiveFunction func(Env, ast.Node, []ast.Node) ast.Node

type Primitive struct {
	Name     string
	Value    primitiveFunction
	MinArity int
	MaxArity int
}

func NewPrimitive(name string, minArity int, maxArity int, value primitiveFunction) *Primitive {
	return &Primitive{
		Name:     name,
		Value:    value,
		MinArity: minArity,
		MaxArity: maxArity,
	}
}

func (p *Primitive) String() string {
	return "#primitive<" + p.Name + ">"
}

func (p *Primitive) isExpr() bool         { return true }
func (p *Primitive) TypeName() string     { return "primitive" }
func (p *Primitive) Loc() *token.Location { return nil }
func (p *Primitive) Equals(n ast.Node) bool {
	panicEvalError(n, "Cannot compare the values of primitive procedures: "+
		p.String()+" and "+n.String())
	return false
}

////////// Function

type Function struct {
	Name       string
	Parameters []ast.Node
	Body       ast.Node
	ParentEnv  Env
	IsMacro    bool
}

func (f *Function) String() string {
	if f.IsMacro {
		return "#macrofunction<" + f.Name + ">"
	}
	return "#function<" + f.Name + ">"
}

func (f *Function) isExpr() bool         { return true }
func (f *Function) Loc() *token.Location { return nil }
func (f *Function) TypeName() string {
	if f.IsMacro {
		return "macrofunction"
	}
	return "function"
}
func (f *Function) Equals(n ast.Node) bool {
	panicEvalError(n, "Cannot compare the values of functions: "+
		f.String()+" and "+n.String())
	return false
}

////////// Macro

type Macro struct {
	Name       string
	Parameters []ast.Node
	Body       ast.Node
	ParentEnv  Env
}

func (m *Macro) String() string {
	return "#macro<" + m.Name + ">"
}

func (m *Macro) isExpr() bool         { return true }
func (m *Macro) TypeName() string     { return "macro" }
func (m *Macro) Loc() *token.Location { return nil }
func (m *Macro) Equals(n ast.Node) bool {
	panicEvalError(n, "Cannot compare the values of macros: "+
		m.String()+" and "+n.String())
	return false
}

////////// Chan

var channelNumber int

type Chan struct {
	id    int
	Value chan ast.Node
}

func NewChan() *Chan {
	cn := channelNumber
	channelNumber++

	return &Chan{
		id:    cn,
		Value: make(chan ast.Node),
	}
}

func (c *Chan) String() string       { return fmt.Sprintf("#chan<%v>", c.id) }
func (c *Chan) isExpr() bool         { return true }
func (c *Chan) Loc() *token.Location { return nil }
func (c *Chan) TypeName() string     { return "chan" }
func (c *Chan) Equals(n ast.Node) bool {
	panicEvalError(n, "Cannot compare the values of chans: "+
		c.String()+" and "+n.String())
	return false
}

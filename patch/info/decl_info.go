package info

import (
	"cmd/compile/internal/syntax"
	xgo_func_name "cmd/compile/internal/xgo_rewrite_internal/patch/func_name"
	"fmt"
)

type DeclKind int

const (
	Kind_Func   DeclKind = 0
	Kind_Var    DeclKind = 1
	Kind_VarPtr DeclKind = 2
	Kind_Const  DeclKind = 3

	// TODO
	// Kind_Interface VarKind = 4
)

func (c DeclKind) IsFunc() bool {
	return c == Kind_Func
}

func (c DeclKind) IsVarOrConst() bool {
	return c == Kind_Var || c == Kind_VarPtr || c == Kind_Const
}
func (c DeclKind) String() string {
	switch c {
	case Kind_Func:
		return "func"
	case Kind_Var:
		return "var"
	case Kind_VarPtr:
		return "var_ptr"
	case Kind_Const:
		return "const"
	default:
		return fmt.Sprintf("decl_%d", int(c))
	}
}

type DeclInfo struct {
	FuncDecl *syntax.FuncDecl
	VarDecl  *syntax.VarDecl

	// when kind == CONST
	ConstDecl *syntax.ConstDecl

	// is this var decl follow a const __xgo_trap_xxx = 1?
	FollowingTrapConst bool

	Kind         DeclKind
	Name         string
	RecvTypeName string
	RecvPtr      bool
	Generic      bool
	Closure      bool
	Stdlib       bool

	// this is an interface type declare
	// only the RecvTypeName is valid
	Interface bool

	// arg names
	RecvName     string
	ArgNames     []string
	ResNames     []string
	FirstArgCtx  bool
	LastResError bool

	FileSyntax *syntax.File
	FileIndex  int

	// this is file name after applied -trimpath
	File string
	Line int
}

func (c *DeclInfo) RefName() string {
	if c.Interface {
		return "nil"
	}
	// if c.Generic, then the ref name is for generic
	if !c.Kind.IsFunc() {
		return c.Name
	}
	return xgo_func_name.FormatFuncRefName(c.RecvTypeName, c.RecvPtr, c.Name)
}

func (c *DeclInfo) RefNameSyntax(pos syntax.Pos) syntax.Expr {
	if c.Interface {
		return syntax.NewName(pos, "nil")
	}
	// if c.Generic, then the ref name is for generic
	if !c.Kind.IsFunc() {
		return syntax.NewName(pos, c.Name)
	}
	return xgo_func_name.FormatFuncRefNameSyntax(pos, c.RecvTypeName, c.RecvPtr, c.Name)
}

func (c *DeclInfo) GenericName() string {
	if !c.Generic {
		return ""
	}
	return c.RefName()
}

func (c *DeclInfo) IdentityName() string {
	if c.Interface {
		return c.RecvTypeName
	}
	if !c.Kind.IsFunc() {
		if c.Kind == Kind_VarPtr {
			return "*" + c.Name
		}
		return c.Name
	}
	return xgo_func_name.FormatFuncRefName(c.RecvTypeName, c.RecvPtr, c.Name)
}

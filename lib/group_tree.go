package lib

import (
	"go/ast"
	"sort"
)

// a tuple for middleware function and weight
type MWTuple struct {
	M *MiddlewareFunc
	W int
}

type GroupTree struct {
	Path     string // /c
	AbsPath  string // /aaa/bbb/c
	Ctx      *GenContext
	Parent   *GroupTree
	Middles  []MWTuple
	Handlers []*HandlerFunc
	Children map[string]*GroupTree
}

func NewGroupTreeNode(path, absPath string, ctx *GenContext, parent *GroupTree) *GroupTree {
	return &GroupTree{
		Path:     path,
		AbsPath:  absPath,
		Ctx:      ctx,
		Parent:   parent,
		Middles:  []MWTuple{},
		Handlers: []*HandlerFunc{},
		Children: make(map[string]*GroupTree),
	}
}

func (node *GroupTree) GenAST() []ast.Stmt {
	root := false
	if node.Parent == nil {
		root = true
		node.Path = EngineName
	}
	var stmts []ast.Stmt
	if root {
		for _, v := range node.Handlers {
			stmts = append(stmts, v.GenAST())
		}
		// ensure the statements in order
		var sortedList []string
		for k, _ := range node.Children {
			sortedList = append(sortedList, k)
		}
		sort.Strings(sortedList)
		for _, v := range sortedList {
			var list = node.Children[v].GenAST()
			stmts = append(stmts, list...)
		}
	} else {
		assign := GenAssign(node.Path, node.Parent.Path, &node.Middles)
		stmts = append(stmts, assign)
		var blockStmt ast.BlockStmt
		for _, v := range node.Handlers {
			blockStmt.List = append(blockStmt.List, v.GenAST())
		}
		for _, v := range node.Children {
			var list = v.GenAST()
			blockStmt.List = append(blockStmt.List, list...)
		}
		stmts = append(stmts, &blockStmt)
	}
	return stmts
}

package lib

import (
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"strings"
)

/*
	the root node does not have ParentGroup , GroupVarName or GroupPath
*/
type GenContext struct {
	Root       *GroupTree
	MiddlePool map[string]*MiddlewareFunc
	GroupPool  map[string]*GroupTree
}

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

func Generate(dirList ...string) {
	var root = new(FileTreeNode)
	for _, dir := range dirList {
		tree := loadFileTree(dir)
		root.Children = append(root.Children, tree)
	}
	var RawFuncList []RawFunc
	WalkFileTree(root, &RawFuncList)

	var hdrs []HandlerFunc
	var mids = make(map[string]*MiddlewareFunc)
	var imports = make(map[string]bool)
	for _, f := range RawFuncList {
		switch x := f.ParseRawFunc().(type) {
		case HandlerFunc:
			hdrs = append(hdrs, x)
			imports[x.RawFunc.PackagePath] = true
		case MiddlewareFunc:
			mids[x.ID] = &x
			imports[x.RawFunc.PackagePath] = true
		}
	}

	var ctx GenContext
	ctx.InitTree(mids, &hdrs)
	f := GenTemplate()
	for k, _ := range imports {
		AddImport(f, k)
	}
	list := ctx.Root.GenAST()
	fileExprListAppend(f, &list)
	export(f)
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
		for _, v := range node.Children {
			var list = v.GenAST()
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

func (f *HandlerFunc) GenAST() *ast.ExprStmt {
	argsList := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + f.RelativePath + "\"",
		},
	}
	// middlewares
	for _, v := range f.Middles {
		argsList = append(argsList, &ast.SelectorExpr{
			X: &ast.Ident{
				Name: v.Package,
			},
			Sel: &ast.Ident{
				Name: v.Signature,
			},
		})
	}
	// handlers
	argsList = append(argsList,
		&ast.SelectorExpr{
			X: &ast.Ident{
				Name: f.Package,
			},
			Sel: &ast.Ident{
				Name: f.Signature,
			},
		})
	stmt := ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: path2VarName(f.Group.Path), // the group where handler in
				},
				Sel: &ast.Ident{
					Name: strings.ToUpper(f.Method),
				},
			},
			Args: argsList,
		},
	}
	return &stmt
}

func (ctx *GenContext) InitTree(mids map[string]*MiddlewareFunc, hdrs *[]HandlerFunc) {
	ctx.MiddlePool = mids
	ctx.GroupPool = make(map[string]*GroupTree)
	ctx.Root = new(GroupTree)
	ctx.Root.Children = make(map[string]*GroupTree)
	// append group handlers
	for i, _ := range *hdrs {
		ctx.PutHandler(&(*hdrs)[i])
	}
	//append group middleware
	for k, _ := range mids {
		ctx.PutMiddleware(mids[k])
	}

}

func (ctx *GenContext) PutMiddleware(m *MiddlewareFunc) {
	for p, w := range m.Groups {
		grp := ctx.GroupPool[p]
		if grp == nil {
			log.Fatalf("no such group:(%s)", p)
		}
		if len(grp.Middles) == 0 {
			grp.Middles = append(grp.Middles, MWTuple{M: m, W: w})
		}else{
			for i, v := range grp.Middles {
				if v.W == w {
					log.Fatalf("same weight for two middlewares: [1]:%s [2]:%s", m.PackagePath, v.M.PackagePath)
				}
				if v.W < w || i == len(grp.Middles)-1 {
					//right of i
					if i != len(grp.Middles)-1 {
						grp.Middles = append(grp.Middles, grp.Middles[i+1:]...)
					}else{
						grp.Middles = append(grp.Middles[:i+1], MWTuple{M: m, W: w})
					}
				}
			}
		}
	}
}

func (ctx *GenContext) PutHandler(h *HandlerFunc) {
	//middlewares
	for _, id := range h.Need {
		if m, ok := ctx.MiddlePool[id]; ok {
			h.Middles = append(h.Middles, m)
		} else {
			log.Fatal("middleware not exist: ", m)
		}
	}
	// grow group tree
	var grp = ctx.Root
	var absPath = "" // absolute path
	if len(h.GroupArray) == 0 {
		grp.Handlers = append(grp.Handlers, h)
	} else {
		// find one's group
		for _, step := range h.GroupArray {
			var find = false
			if _, ok := grp.Children[step]; ok {
				find = true
			}
			absPath += step
			if !find {
				node := NewGroupTreeNode(step, absPath, ctx, grp)
				grp.Children[step] = node
				ctx.GroupPool[absPath] = node
			}
			grp = grp.Children[step]
		}
		grp.Handlers = append(grp.Handlers, h)
	}
	h.Group = grp
}

// /var/:name  => varName
func path2VarName(path string) string {
	ss := strings.Split(path, "/")
	if ss[0] == "" {
		ss = ss[1:]
	}
	for i, _ := range ss {
		ss[i] = strings.Trim(ss[i], ":*")
		ss[i] = strings.ToLower(ss[i])
	}
	var name = ss[0]
	for i := 1; i < len(ss); i++ {
		name += strings.ToUpper(ss[i][:1]) + ss[i][1:]
	}
	return name
}

func GenAssign(path, pPath string, mids *[]MWTuple) *ast.AssignStmt {
	args := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + path + "\"",
		},
	}
	for _, v := range *mids {
		args = append(args, &ast.SelectorExpr{
			X: &ast.Ident{
				Name: v.M.Package,
			},
			Sel: &ast.Ident{
				Name: v.M.Signature,
			},
		})
	}
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: path2VarName(path),
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: path2VarName(pPath),
					},
					Sel: &ast.Ident{
						Name: "Group",
					},
				},
				Args: args,
			},
		},
	}
}

func export(f *ast.File) {
	set := token.NewFileSet()
	file, _ := os.Create(DefaultFileName)
	_, _ = file.WriteString(Tips)
	err := format.Node(file, set, f)
	if err != nil {
		panic("can not export file")
	}
}

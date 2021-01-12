package lib

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

const (
	EngineName      = "e"
	Tips            = `//======================================================================//
//																		//
//		Generated by gin-annotation (github.com/1-st/gin-annotation)	//
//																		//
//				***	you should call Route() manually ***				//
//																		//
//======================================================================//


`
	TemplateFile = `
	package _example

	import (
		"github.com/gin-gonic/gin"
	)

	func Route(e *gin.Engine) {
	}`
)

var (
	DefaultFileName = "route.entry.go"
)

func GenTemplate() *ast.File {
	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, "", TemplateFile, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	p, _ := parser.ParseDir(fSet, "./", nil, parser.ParseComments)
	var curPkgName string
	for curPkgName, _ = range p {
		break
	}
	f.Name.Name = curPkgName
	return f
}

func AddImport(f *ast.File, im string) {
	f.Decls[0].(*ast.GenDecl).Specs = append(f.Decls[0].(*ast.GenDecl).Specs, &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + im + "\"",
		},
	})
}

func fileExprListAppend(f *ast.File, list *[]ast.Stmt) {
	if len(f.Decls) != 2 {
		log.Fatal("invalid file")
	}
	f.Decls[1].(*ast.FuncDecl).Body.List = *list
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


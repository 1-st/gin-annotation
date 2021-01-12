package lib

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strconv"
	"strings"
)

type Func interface {
	ParseRawFunc() Func
}

type HandlerFunc struct {
	RawFunc
	Group        *GroupTree
	GroupArray   []string // /user -> /name
	Need         []string // middleware id
	Middles      []*MiddlewareFunc
	Method       string
	RelativePath string
}

type MiddlewareFunc struct {
	RawFunc
	ID    string
	Group map[string]int
}

type RawFunc struct {
	Comment     string
	PackagePath string
	Package     string
	Signature   string
}

// ParseRawFunc generate Func from RawFunc
func (f RawFunc) ParseRawFunc() Func {
	var (
		leftBracket  = -1
		rightBracket = -1
	)
	for i := 0; i < len(f.Comment); i++ {
		if f.Comment[i] == '[' {
			leftBracket = i
		} else if f.Comment[i] == ']' {
			rightBracket = i
			break
		}
	}
	if leftBracket == -1 || rightBracket == -1 {
		panic("error: no [] found")
		return nil
	}
	comment := f.Comment[leftBracket+1 : rightBracket]
	ss := strings.Split(comment, ",")
	for i, v := range ss {
		ss[i] = strings.Trim(v, "\t\n")
	}
	var attr = make(map[string][]string)
	for _, v := range ss {
		var kv []string
		for i := 0; i < len(v); i++ {
			if v[i] == ':' {
				kv = append(kv, v[:i])
				kv = append(kv, v[i+1:])
				break
			}
		}
		if len(kv) != 2 {
			panic(fmt.Sprintf("error kv : %v", kv))
		}
		var key = strings.ToLower(kv[0])
		var value = strings.ToLower(kv[1])
		if _, ok := attr[key]; !ok {
			attr[key] = make([]string, 0)
		}
		attr[key] = append(attr[key], value)
	}
	var checkExist = func(m map[string][]string, attrs ...string) bool {
		for _, v := range attrs {
			if _, ok := m[v]; !ok {
				return false
			}
		}
		return true
	}
	matchHandler := checkExist(attr, "method", "path")
	matchMiddleware := checkExist(attr, "id")
	if matchHandler &&
		!matchMiddleware {
		var groups []string
		var need []string
		if v, ok := attr["groups"]; ok {
			groups = strings.Split(v[0], " ")
		}
		if v, ok := attr["need"]; ok {
			need = strings.Split(v[0], " ")
		}
		var h HandlerFunc
		h.RawFunc = f
		h.GroupArray = groups
		h.Method = formatMethod(attr["method"][0])
		h.RelativePath = attr["path"][0]
		h.Need = need
		return h
	} else if matchMiddleware &&
		!matchHandler {
		var m MiddlewareFunc
		m.RawFunc = f
		m.Group = make(map[string]int)
		m.ID = attr["id"][0]
		for _, v := range attr["group"] {
			ss := strings.Split(v, "@")
			n, err := strconv.ParseInt(ss[1], 10, 64)
			if err != nil {
				log.Print("error: sortIndex is not integer", err)
				return nil
			}
			m.Group[ss[0]] = int(n)
		}
		return m
	} else {
		panic("no match FuncType")
	}
}

// TraverseRawFuncList traverses list of RawFunc and get handler list ,middleware map and imports list
func TraverseRawFuncList(list *[]RawFunc, h *[]HandlerFunc, m map[string]*MiddlewareFunc, i map[string]bool) {
	for _, f := range *list {
		switch x := f.ParseRawFunc().(type) {
		case HandlerFunc:
			*h = append(*h, x)
			i[x.RawFunc.PackagePath] = true
		case MiddlewareFunc:
			m[x.ID] = &x
			i[x.RawFunc.PackagePath] = true
		}
	}
}

func (h *HandlerFunc) GenAST() *ast.ExprStmt {
	argsList := []ast.Expr{
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + h.RelativePath + "\"",
		},
	}
	// middlewares
	for _, v := range h.Middles {
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
				Name: h.Package,
			},
			Sel: &ast.Ident{
				Name: h.Signature,
			},
		})
	stmt := ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: path2VarName(h.Group.Path), // the group where handler in
				},
				Sel: &ast.Ident{
					Name: h.Method,
				},
			},
			Args: argsList,
		},
	}
	return &stmt
}

func funcIsController(f *ast.FuncDecl) bool {
	if f.Type.Results != nil {
		return false
	}
	if len(f.Type.Params.List) != 1 {
		return false
	}
	starExpr, ok := f.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	if ident.Name != "gin" {
		return false
	}
	if selectorExpr.Sel.Name != "Context" {
		return false
	}
	return true
}

func formatMethod(str string) string {
	low := strings.ToLower(str)
	switch low {
	case "get":
		return "GET"
	case "post":
		return "POST"
	case "put":
		return "PUT"
	case "delete":
		return "DELETE"
	case "options":
		return "OPTIONS"
	case "patch":
		return "PATCH"
	case "any":
		return "Any"
	default:
		panic("invalid method-annotation: " + str)
	}
}

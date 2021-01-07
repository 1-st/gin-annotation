package lib

import (
	"go/ast"
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
	ID     string
	Groups map[string]int
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
		log.Print("error: no [] found")
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
			log.Print("error kv : ", kv)
			return nil
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
		var group []string
		var need []string
		if v, ok := attr["group"]; ok {
			group = strings.Split(v[0], " ")
		}
		if v, ok := attr["need"]; ok {
			need = strings.Split(v[0], " ")
		}
		var h HandlerFunc
		h.RawFunc = f
		h.GroupArray = group
		h.Method = attr["method"][0]
		h.RelativePath = attr["path"][0]
		h.Need = need
		return h
	} else if matchMiddleware &&
		!matchHandler {
		var m MiddlewareFunc
		m.RawFunc = f
		m.Groups = make(map[string]int)
		m.ID = attr["id"][0]
		for _, v := range attr["groups"] {
			ss := strings.Split(v, "@")
			n, err := strconv.ParseInt(ss[1], 10, 64)
			if err != nil {
				log.Print("error: sortIndex is not integer", err)
				return nil
			}
			m.Groups[ss[0]] = int(n)
		}
		return m
	} else {
		log.Print("error: no match FuncType")
		return nil
	}
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

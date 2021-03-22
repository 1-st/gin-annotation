package lib

import "strings"

// entry function
// Generate the entry function of project
func Generate(dirList ...string) {
	var root = NewEmptyFileTree()
	for _, dir := range dirList {
		tree := LoadFileTree(dir)
		root.Children = append(root.Children, tree)
	}

	var RawFuncList []RawFunc
	WalkFileTree(root, &RawFuncList)

	var handlers []HandlerFunc
	var middlewares = make(map[string]*MiddlewareFunc)
	var imports = make(map[string]bool)
	TraverseRawFuncList(&RawFuncList, &handlers, middlewares, imports)

	var ctx GenContext
	ctx.InitGroupTree(middlewares, &handlers)

	tplFile := GenTemplate()
	for k, _ := range imports {
		AddImport(tplFile, k)
	}

	list := ctx.Root.GenAST()
	fileExprListAppend(tplFile, &list)

	export(tplFile)
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

package lib

import (
	"gin-annotation/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"path/filepath"
)

type FileTreeNode struct {
	PackagePath string
	Files       map[string]*ast.File
	Children    []*FileTreeNode
}

func loadFileTree(dir string) *FileTreeNode {
	var node FileTreeNode
	dfsLoadFileTree(&node, dir)
	return &node
}

func dfsLoadFileTree(node *FileTreeNode, dir string) {
	p, _ := packages.Load(&packages.Config{Dir: dir})
	node.PackagePath = p[0].PkgPath
	var fSet token.FileSet
	absDirPath, _ := filepath.Abs(dir)
	parsedPackages, _ := parser.ParseDir(&fSet, absDirPath, nil, parser.ParseComments)
	for _, v := range parsedPackages {
		node.Files = v.Files
	}
	dirs := utils.GetChildDirs(dir)
	if len(dirs) != 0 {
		for _, v := range dirs {
			var newNode FileTreeNode
			node.Children = append(node.Children, &newNode)
			dfsLoadFileTree(&newNode, dir+"/"+v)
		}
	}
}

// WalkFileTree walk file tree and get the raw function
func WalkFileTree(node *FileTreeNode, list *[]RawFunc) {
	for _, v := range node.Files {
		ast.Inspect(v, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if funcIsController(x) {
					*list = append(*list, RawFunc{
						Comment:     x.Doc.Text(),
						PackagePath: node.PackagePath,
						Signature:   x.Name.Name,
						Package: v.Name.Name,
					})
				}
			}
			return true
		})
	}
	if node.Children != nil {
		for _, v := range node.Children {
			WalkFileTree(v, list)
		}
	}
}



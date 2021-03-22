package lib

import (
	"fmt"
	"github.com/1-st/gin-annotation/utils"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
)

//load files to a tree data structure
type FileTreeNode struct {
	PackagePath string
	Files       map[string]*ast.File
	Children    []*FileTreeNode
}

func NewEmptyFileTree() *FileTreeNode {
	return &FileTreeNode{
		PackagePath: "",
		Files:       nil,
		Children:    nil,
	}
}

// LoadFileTree load disk file into a tree
func LoadFileTree(dir string) *FileTreeNode {
	var node FileTreeNode
	dfsLoadFileTree(&node, dir)
	return &node
}

// WalkFileTree walk file tree and get the raw function
func WalkFileTree(from *FileTreeNode, to *[]RawFunc) {
	for _, v := range from.Files {
		ast.Inspect(v, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if funcIsController(x) {
					r := &RawFunc{
						Comment:     x.Doc.Text(),
						PackagePath: from.PackagePath,
						Package:     v.Name.Name,
						Signature:   x.Name.Name,
					}
					logFunction(r)
					*to = append(*to, *r)
				}
			}
			return true
		})
	}
	if from.Children != nil {
		for _, v := range from.Children {
			WalkFileTree(v, to)
		}
	}
}

func logFunction(f *RawFunc) {
	fmt.Println("detected: " + f.Package + "." + f.Signature + " in " + f.PackagePath)
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

func export(f *ast.File) {
	set := token.NewFileSet()
	file, _ := os.Create(DefaultFileName)
	_, _ = file.WriteString(Tips)
	err := format.Node(file, set, f)
	if err != nil {
		panic("can not export file")
	}
}

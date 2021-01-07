package main

import (
	"fmt"
	"gin-annotation/lib"
	"gin-annotation/utils"
	"os"
)

func printHelp() {
	fmt.Print("_example:\n")
	fmt.Print("gin-annotation ./dir1 ./dir2 \n")
}

func main() {
	if os.Getenv("GIN_ANNOTATION_FILE") != "" {
		lib.DefaultFileName = os.Getenv("GIN_ANNOTATION_FILE")
	}
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	dirList := os.Args[1:]
	var checkDir = func(dir string) {
		if !utils.IsDir(dir) {
			fmt.Println(dir + " is not _example dir")
			os.Exit(1)
		}
	}
	for _, v := range dirList {
		checkDir(v)
	}
	lib.Generate(dirList...)
}

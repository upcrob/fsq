package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: test <expression>")
		os.Exit(1)
	}

	lexer := new(Lexer)
	lexer.expr = os.Args[1]
	yyParse(lexer)

	if programRoot == nil {
		fmt.Println("invalid expression")
		os.Exit(1)
	}

	if !validAttributesRequested() {
		fmt.Println("invalid expression - invalid attributes requested")
		os.Exit(1)
	}

	// parse tree optimizations
	optimizeContentsContains(programRoot.children[3])

	// walk file system
	filepath.Walk(programRoot.children[2].sval, eval)
}

func eval(path string, file os.FileInfo, err error) error {
	if file == nil {
		return nil
	}

	if evaluate(path, file, programRoot.children[3]) {
		printRelevant(path, file)
	}
	return nil
}

func printRelevant(path string, file os.FileInfo) {
	printStarted := false
	if attributeRequested(T_PATH) {
		if file.IsDir() {
			path += "/"
		}
		fmt.Print(path)
		printStarted = true
	}
	if attributeRequested(T_NAME) {
		if printStarted {
			fmt.Print(",")
		}
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		fmt.Print(name)
		printStarted = true
	}
	if attributeRequested(T_SIZE) {
		if printStarted {
			fmt.Print(",")
		}
		fmt.Print(int(file.Size() / 1048576))
	}
	fmt.Println()
}

func attributeRequested(ntype int) bool {
	attribs := programRoot.children[0].children
	for _, v := range attribs {
		if v.ntype == ntype {
			return true
		}
	}
	return false
}

func validAttributesRequested() bool {
	attribs := programRoot.children[0].children
	for _, v := range attribs {
		if !(v.ntype == T_NAME || v.ntype == T_PATH || v.ntype == T_SIZE) {
			return false
		}
	}
	return true
}

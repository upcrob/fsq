package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const TIMESTAMP_FORMAT = "01/02/2006 15:04:05"
const DATE_FORMAT = "01/02/2006"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: fsq <expression>")
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
	shiftExpressionRight(programRoot.children[3], isContentStartswithExpression)
	shiftExpressionRight(programRoot.children[3], isContentEndswithExpression)
	shiftExpressionRight(programRoot.children[3], isContentContainsExpression)

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
	if attributeRequested(T_MODIFIED) {
		fmt.Print(file.ModTime().Format(TIMESTAMP_FORMAT) + "  ")
	}
	if attributeRequested(T_SIZE) {
		fmt.Print(pad(strconv.Itoa(int(file.Size())), 11) + " ")
	}
	if attributeRequested(T_NAME) {
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		fmt.Print(pad(name, 25) + " ")
	}
	if attributeRequested(T_PATH) {
		if file.IsDir() {
			path += "/"
		}
		fmt.Print(path)
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
		if !(v.ntype == T_NAME || v.ntype == T_PATH || v.ntype == T_SIZE || v.ntype == T_MODIFIED) {
			return false
		}
	}
	return true
}

func pad(str string, size int) string {
	ilen := len(str)
	for i := 0; i < size - ilen; i++ {
		str += " "
	}
	return str
}

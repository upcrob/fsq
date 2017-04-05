package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const TIMESTAMP_FORMAT = "01/02/2006 15:04:05"
const TIMESTAMP_FORMAT_PATTERN = "MM/DD/YYYY hh:mm:ss"
const DATE_FORMAT = "01/02/2006"
const DATE_FORMAT_PATTERN = "MM/DD/YYYY"
const VERSION = "1.4.0"

type result struct {
	order   int
	matched bool
	path    string
	file    os.FileInfo
}

var evalGroup sync.WaitGroup
var printGroup sync.WaitGroup
var resultChannel chan *result
var count int
var searchStrings []SearchString

func main() {
	if len(os.Args) < 2 {
		fmt.Println("fsq version " + VERSION)
		fmt.Println("usage: fsq <expression>")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		execute_expression(os.Args[1])
		return
	}

	buf := &bytes.Buffer{}
	wr := bufio.NewWriter(buf)
	for i := 1; i < len(os.Args); i++ {
		wr.WriteString(os.Args[i])
		wr.WriteString(" ")
	}
	wr.Flush()
	execute_expression(buf.String())
}

func execute_expression(expr string) {
	count = 0
	resultChannel = make(chan *result)

	lexer := new(Lexer)
	lexer.expr = expr
	yyParse(lexer)

	if programRoot == nil {
		fmt.Println("invalid expression")
		os.Exit(1)
	}

	if !validAttributesRequested() {
		fmt.Println("invalid expression - invalid attributes requested")
		os.Exit(1)
	}

	// validate parse tree
	validateTree(programRoot.children[3])

	// parse tree optimization
	shiftShortestPathLeft(programRoot.children[3])

	// collect any strings to search files for preemptively
	searchStrings = collectFileSearchStrings(programRoot.children[3])

	// start print routine, this will print out the results sent from
	// doEval() via resultChannel
	go printRoutine()
	printGroup.Add(1)

	// walk file system
	// the grammar causes reverses the list of root directories,
	// so we have to walk the list in reverse
	rootList := programRoot.children[2].children
	for i := len(rootList) - 1; i >= 0; i-- {
		root := rootList[i].sval
		filepath.Walk(root, eval)
	}
	evalGroup.Wait()
	resultChannel <- nil
	printGroup.Wait()
}

func eval(path string, file os.FileInfo, err error) error {
	if file == nil {
		return nil
	}

	rootList := programRoot.children[2].children
	for i := 0; i < len(rootList); i++ {
		if path == rootList[i].sval {
			// exclude root directory
			return nil
		}
	}

	evalGroup.Add(1)
	go doEval(path, file, programRoot.children[3], count)
	count++
	return nil
}

func doEval(path string, file os.FileInfo, n *tnode, order int) {
	fileSearch := newFileSearch(searchStrings, path)

	res := new(result)
	res.order = order
	normalizedPath := forwardSlashes(path)
	if evaluate(normalizedPath, file, n, fileSearch) {
		res.matched = true
		res.path = normalizedPath
		res.file = file
	} else {
		res.matched = false
	}
	resultChannel <- res
	evalGroup.Done()
}

func printRoutine() {
	current := 0
	cache := make(map[int]*result)
	for {
		res := <-resultChannel
		if res == nil {
			break
		}

		if res.order == current {
			if res.matched {
				// this is the next item to print out, print it
				printRelevant(res.path, res.file)
			}
			current++

			// print subsequent items available in the cache
			cached := cache[current]
			for cached != nil {
				if cached.matched {
					printRelevant(cached.path, cached.file)
				}
				delete(cache, current)
				current++
				cached = cache[current]
			}
		} else {
			// store the item in the cache
			cache[res.order] = res
		}
	}
	printGroup.Done()
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
	for i := 0; i < size-ilen; i++ {
		str += " "
	}
	return str
}

func forwardSlashes(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

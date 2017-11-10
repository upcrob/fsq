package main

import (
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
const VERSION = "1.7.1"

type result struct {
	order   int
	matched bool
	path    string
	file    os.FileInfo
}

var evalGroup sync.WaitGroup
var printGroup sync.WaitGroup
var countMutex sync.Mutex
var resultChannel chan *result
var count int
var searchStrings []SearchString
var matchCount int
var sizeCount int64

func main() {
	if len(os.Args) != 2 {
		fmt.Println("fsq version " + VERSION)
		fmt.Println("usage: fsq <expression>")
		os.Exit(1)
	}
	execute_expression(os.Args[1])
}

func execute_expression(expr string) {
	count = 0
	matchCount = 0
	sizeCount = 0
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
	reg := compileRegexes(programRoot.children[3])
	if reg != "" {
		fmt.Println("invalid regex: " + reg)
		os.Exit(1)
	}

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

	if attributeRequested(T_STATS) {
		// print aggregate statistics
		attribs := programRoot.children[0].children
		if len(attribs) > 1 && matchCount > 0 {
			// non-statistics output was printed, add a newline
			fmt.Println()
		}
		fmt.Println("files: " + strconv.Itoa(matchCount))

		fsize := friendlySize(sizeCount)
		_, err := strconv.Atoi(fsize)
		if err != nil {
			fsize = strconv.FormatInt(sizeCount, 10) + " (" + fsize + ")"
		}
		fmt.Println("size: " + fsize)
	}
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
		countMutex.Lock()
		matchCount++
		sizeCount += file.Size()
		countMutex.Unlock()
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
	anyRequested := false
	if attributeRequested(T_MODIFIED) {
		anyRequested = true
		fmt.Print(file.ModTime().Format(TIMESTAMP_FORMAT) + "  ")
	}
	if attributeRequested(T_SIZE) {
		anyRequested = true
		fmt.Print(pad(strconv.Itoa(int(file.Size())), 11) + " ")
	}
	if attributeRequested(T_FSIZE) {
		anyRequested = true
		fmt.Print(pad(friendlySize(file.Size()), 6) + " ")
	}
	if attributeRequested(T_NAME) {
		anyRequested = true
		name := file.Name()
		if file.IsDir() {
			name += "/"
		}
		fmt.Print(pad(name, 25) + " ")
	}
	if attributeRequested(T_PATH) {
		anyRequested = true
		if file.IsDir() {
			path += "/"
		}
		fmt.Print(path)
	}

	if anyRequested {
		fmt.Println()
	}
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
		if !(v.ntype == T_NAME || v.ntype == T_PATH || v.ntype == T_SIZE ||
			v.ntype == T_MODIFIED || v.ntype == T_STATS || v.ntype == T_FSIZE) {
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

func friendlySize(bytes int64) string {
	if bytes <= 1000 {
		return strconv.Itoa(int(bytes))
	}
	fbytes := float64(bytes) / 1000.0
	if fbytes <= 1000.0 {
		return strconv.FormatFloat(fbytes, 'f', 1, 64) + "k"
	}
	fbytes /= 1000.0
	if fbytes <= 1000.0 {
		return strconv.FormatFloat(fbytes, 'f', 1, 64) + "m"
	}
	fbytes /= 1000.0
	return strconv.FormatFloat(fbytes, 'f', 1, 64) + "g"
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"runtime"
	"flag"
	"io/ioutil"
)

const TIMESTAMP_FORMAT = "01/02/2006 15:04:05"
const TIMESTAMP_FORMAT_PATTERN = "MM/DD/YYYY hh:mm:ss"
const DATE_FORMAT = "01/02/2006"
const DATE_FORMAT_PATTERN = "MM/DD/YYYY"
const VERSION = "1.9.0"

type result struct {
	order   int
	matched bool
	path    string
	file    os.FileInfo
	hash    *ComputedHash
}

var evalGroup sync.WaitGroup
var printGroup sync.WaitGroup
var countMutex sync.Mutex
var resultChannel chan *result
var routineTickets chan int
var count int
var searchStrings []SearchString
var matchCount int
var sizeCount int64
var maxRoutines int
var doDelete bool

func main() {
	flag.BoolVar(&doDelete, "d", false, "delete matched files")
	flag.Parse()

	if len(flag.Args()) <= 0 {
		fmt.Println("fsq version " + VERSION)
		fmt.Println("usage: fsq [-d] <expression>")
		fmt.Println("  -d Delete matched files and directories")
		os.Exit(1)
	}

	maxRoutines = runtime.NumCPU()
	execute_expression(flag.Arg(0))
}

func includesPath(paths []*tnode, path string) bool {
	if paths == nil {
		return false
	}

	for i := 0; i < len(paths); i++ {
		if paths[i].sval == path {
			return true
		}
	}
	return false
}

func walk(path string, eval func(path string, file os.FileInfo), excludePaths []*tnode) {
	if !includesPath(excludePaths, path) {
		files, err := ioutil.ReadDir(path)
		if err == nil {
			for i := 0; i < len(files); i++ {
				var newPath string
				if path == "." {
					newPath = files[i].Name()
				} else {
					newPath = path + "/" + files[i].Name()
				}
				if files[i].IsDir() {
					eval(newPath, files[i])
					walk(newPath, eval, excludePaths)
				} else {
					eval(newPath, files[i])
				}
			}
		}
	}
}

func execute_expression(expr string) {
	count = 0
	matchCount = 0
	sizeCount = 0
	resultChannel = make(chan *result)

	// initialize the pool with the maximum number of tickets (goroutines)
	// that may be running at a given time
	routineTickets = make(chan int, maxRoutines)
	for i := 0; i < maxRoutines; i++ {
		routineTickets <- 1
	}

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

	if programRoot.children[1].ival == 1 || programRoot.children[1].ival == 0 {
		// walk file system
		rootList := programRoot.children[2].children

		// the grammar reverses the list of root directories,
		// so we have to walk the list in reverse
		for i := len(rootList) - 1; i >= 0; i-- {
			root := rootList[i].sval
			walk(root, eval, nil)
		}
	} else {
		// walk the file system, excluding the list of provided paths
		walk(".", eval, programRoot.children[2].children)
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

func eval(path string, file os.FileInfo) {
	if file == nil {
		return
	}

	rootList := programRoot.children[2].children
	for i := 0; i < len(rootList); i++ {
		if path == rootList[i].sval {
			// exclude root directory
			return
		}
	}

	evalGroup.Add(1)
	takeRoutineTicket()
	go doEval(path, file, programRoot.children[3], count)
	count++
}

func doEval(path string, file os.FileInfo, n *tnode, order int) {
	fileSearch := newFileSearch(searchStrings, path)
	computedHash := new(ComputedHash)

	res := new(result)
	res.order = order
	normalizedPath := forwardSlashes(path)
	if evaluate(normalizedPath, file, n, fileSearch, computedHash) {
		res.matched = true
		res.path = normalizedPath
		res.file = file
		res.hash = computedHash
		countMutex.Lock()
		matchCount++
		sizeCount += file.Size()
		countMutex.Unlock()
	} else {
		res.matched = false
	}
	resultChannel <- res
	evalGroup.Done()
	returnRoutineTicket()
}

func handleDelete(path string, isDir bool) bool {
	if (doDelete) {
		if isDir {
			err := os.RemoveAll(path)
			return err == nil
		} else {
			err := os.Remove(path)
			if os.IsNotExist(err) {
				// file was already deleted by either parent directory or external process
				return true
			} else {
				return err == nil
			}
		}
	} else {
		return false
	}
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
				deleted := handleDelete(res.path, res.file.IsDir())
				printRelevant(res.path, res.file, res.hash, deleted)
			}
			current++

			// print subsequent items available in the cache
			cached := cache[current]
			for cached != nil {
				if cached.matched {
					deleted := handleDelete(cached.path, cached.file.IsDir())
					printRelevant(cached.path, cached.file, res.hash, deleted)
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

func printRelevant(path string, file os.FileInfo, hash *ComputedHash, deleted bool) {
	if (doDelete) {
		if (deleted) {
			fmt.Print("(deleted) ")
		} else {
			fmt.Print("(delete failed) ")
		}
	}
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
	if attributeRequested(T_SHA1) {
		anyRequested = true
		var sha string = ""
		if (hash == nil || hash.sha1 == nil) {
			if (file.IsDir()) {
				sha = ""
			} else {
				sha = getFileSha1(path)
			}
		} else {
			sha = *hash.sha1
		}
		fmt.Print(pad(sha, 40) + " ")
	}
	if attributeRequested(T_MD5) {
		anyRequested = true
		var md5 string = ""
		if (hash == nil || hash.md5 == nil) {
			if (file.IsDir()) {
				md5 = ""
			} else {
				md5 = getFileMd5(path)
			}
		} else {
			md5 = *hash.md5
		}
		fmt.Print(pad(md5, 32) + " ")
	}
	if attributeRequested(T_SHA256) {
		anyRequested = true
		var sha256 string = ""
		if (hash == nil || hash.sha256 == nil) {
			if (file.IsDir()) {
				sha256 = ""
			} else {
				sha256 = getFileSha256(path)
			}
		} else {
			sha256 = *hash.sha256
		}
		fmt.Print(pad(sha256, 64) + " ")
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
			v.ntype == T_MODIFIED || v.ntype == T_STATS || v.ntype == T_FSIZE ||
			v.ntype == T_SHA1 || v.ntype == T_MD5 || v.ntype == T_SHA256) {
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

func takeRoutineTicket() {
	<- routineTickets
}

func returnRoutineTicket() {
	routineTickets <- 1
}

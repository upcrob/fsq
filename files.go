package main

import (
	"os"
	"strings"
	"regexp"
	"io/ioutil"
)

const DEFAULT_BLOCK_SIZE = 1024

type SearchString struct {
	str string
	caseSensitive bool
}

type FileSearch struct {
	file *os.File
	path string
	buff1 []byte
	buff1Len int
	buff2 []byte
	buff2Len int
	blockSize int
	blocksRead int
	closed bool
	searchStrings []SearchString
	contains []SearchString
	read bool
	content string
}

func newFileSearch(searchStrings []SearchString, path string) *FileSearch {
	max := 0
	for i := 0; i < len(searchStrings); i++ {
		if len(searchStrings[i].str) > max {
			max = len(searchStrings[i].str)
		}
	}

	fs := new(FileSearch)
	fs.path = path
	fs.closed = false
	fs.searchStrings = searchStrings
	fs.blocksRead = 0
	fs.contains = make([]SearchString, 0, 5)
	fs.read = false
	if max > DEFAULT_BLOCK_SIZE {
		fs.blockSize = max
	} else {
		fs.blockSize = DEFAULT_BLOCK_SIZE
	}
	fs.buff1 = make([]byte, fs.blockSize)
	fs.buff2 = make([]byte, fs.blockSize)
	return fs
}

func updateSearch(fs *FileSearch, info os.FileInfo) {
	if fs.closed {
		// this search has been closed, return
		return
	} else if fs.file == nil {
		// the file has not been opened yet, open it
		var err error
		fs.file, err = os.Open(fs.path)
		if err != nil {
			// something went wrong, mark the search as closed
			fs.closed = true
			return
		}
	}

	// read bytes
	var bytesRead int
	var err error
	// alternate between byte buffers
	if fs.blocksRead % 2 == 1 {
		bytesRead, err = fs.file.Read(fs.buff1)
		fs.buff1Len = bytesRead
	} else {
		bytesRead, err = fs.file.Read(fs.buff2)
		fs.buff2Len = bytesRead
	}
	if err != nil {
		fs.closed = true
		return
	} else if int64(fs.blocksRead * fs.blockSize + bytesRead) == info.Size() {
		defer fs.file.Close()
		fs.closed = true
	}

	var sval string
	if fs.blocksRead % 2 == 0 {
		sval = string(fs.buff1[:fs.buff1Len]) + string(fs.buff2[:fs.buff2Len])
	} else {
		sval = string(fs.buff2[:fs.buff2Len]) + string(fs.buff1[:fs.buff1Len])
	}
	svalLower := strings.ToLower(sval)
	fs.blocksRead++

	// check for file containing target strings
	for i := 0; i < len(fs.searchStrings); i++ {
		if fs.searchStrings[i].caseSensitive {
			if strings.Contains(sval, fs.searchStrings[i].str) {
				fs.contains = append(fs.contains, fs.searchStrings[i])
			}
		} else {
			if strings.Contains(svalLower, strings.ToLower(fs.searchStrings[i].str)) {
				fs.contains = append(fs.contains, fs.searchStrings[i])
			}
		}
	}
}

func fileContainsString(fs *FileSearch, info os.FileInfo, target string, caseSensitive bool) bool {
	if fs.read {
		// file contents had to be fully read into memory previously, use
		// this instead of reading the file again incrementally
		if caseSensitive {
			return strings.Contains(fs.content, target)
		} else {
			return strings.Contains(strings.ToLower(fs.content), strings.ToLower(target))
		}
	}

	for !fs.closed && !searchStringExists(fs.contains, SearchString{target, caseSensitive}) {
		updateSearch(fs, info)
	}
	return searchStringExists(fs.contains, SearchString{target, caseSensitive})
}

func searchStringExists(searchSlice []SearchString, searchString SearchString) bool {
	for i := 0; i < len(searchSlice); i++ {
		if searchSlice[i].caseSensitive == searchString.caseSensitive &&
				searchSlice[i].str == searchString.str {
			return true
		}
	}
	return false
}

func fileMatchesString(fs *FileSearch, info os.FileInfo, re *regexp.Regexp) bool {
	if !fs.read {
		// read whole file into memory
		bytes, err := ioutil.ReadFile(fs.path)
		if err != nil {
			return false
		}
		fs.content = string(bytes)
		fs.read = true
	}

	return re.MatchString(fs.content)
}

func fileStartsWithString(path string, str string, caseSensitive bool) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buff := make([]byte, len(str))
	_, readErr := f.Read(buff)
	if readErr != nil {
		return false
	}

	startValue := string(buff)
	if !caseSensitive {
		startValue = strings.ToLower(startValue)
	}

	return str == startValue
}

func fileEndsWithString(path string, info os.FileInfo, str string, caseSensitive bool) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	bsize := len(str)
	buff := make([]byte, bsize)
	cbuff := make([]byte, 1)

	j := bsize - 1
	for i := info.Size() - 1; i >= 0 && j >= 0; i-- {
		_, readErr := f.ReadAt(cbuff, i)
		if readErr != nil {
			return false
		}

		if cbuff[0] >= 32 || j < bsize-1 {
			buff[j] = cbuff[0]
			j--
		}
	}

	endValue := string(buff)
	if !caseSensitive {
		endValue = strings.ToLower(endValue)
	}

	return str == endValue
}

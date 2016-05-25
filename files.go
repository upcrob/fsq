package main

import (
	"os"
	"strings"
	"unicode"
)

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

func fileContainsString(path string, str string, caseSensitive bool) bool {
	bsize := len(str)
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buff := 1
	buff1 := make([]byte, bsize)
	buff2 := make([]byte, bsize)

	// read in a block
	var readErr error = nil
	for readErr == nil {
		if buff == 1 {
			buff = 2
			_, readErr = f.Read(buff1)

			if !caseSensitive {
				for i := 0; i < bsize; i++ {
					buff1[i] = byte(unicode.ToLower(rune(buff1[i])))
				}
			}

			if strings.Contains(string(string(buff2)+string(buff1)), str) {
				return true
			}
		} else {
			buff = 1
			_, readErr = f.Read(buff2)

			if !caseSensitive {
				for i := 0; i < bsize; i++ {
					buff2[i] = byte(unicode.ToLower(rune(buff2[i])))
				}
			}

			if strings.Contains(string(string(buff1)+string(buff2)), str) {
				return true
			}
		}
	}

	return false
}

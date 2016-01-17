package main

import (
	"os"
	"strings"
)

func fileContainsString(path string, str string) bool {
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

			if strings.Contains(string(string(buff2) + string(buff1)), str) {
				return true
			}
		} else {
			buff = 1
			_, readErr = f.Read(buff2)

			if strings.Contains(string(string(buff1) + string(buff2)), str) {
				return true
			}
		}
	}

	return false
}

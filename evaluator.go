package main

import (
	"os"
	"fmt"
	"time"
	"strings"
)

func evaluate(path string, info os.FileInfo, n *tnode) bool {
	if n.ntype == T_STARTSWITH {
		if strings.HasPrefix(resolveAsString(path, left(n), info), right(n).sval) {
			return true
		}
	} else if n.ntype == T_ENDSWITH {
		if strings.HasSuffix(resolveAsString(path, left(n), info), right(n).sval) {
			return true
		}
	} else if n.ntype == T_GT {
		if resolveAsInt(left(n), info) > resolveAsInt(right(n), info) {
			return true
		}
	} else if n.ntype == T_GTE {
		if resolveAsInt(left(n), info) >= resolveAsInt(right(n), info) {
			return true
		}
	} else if n.ntype == T_LT {
		if resolveAsInt(left(n), info) < resolveAsInt(right(n), info) {
			return true
		}
	} else if n.ntype == T_LTE {
		if resolveAsInt(left(n), info) <= resolveAsInt(right(n), info) {
			return true
		}
	} else if n.ntype == T_EQ {
		if left(n).ntype == T_NAME {
			if resolveAsString(path, left(n), info) == right(n).sval {
				return true
			}
		} else if left(n).ntype == T_PATH {
			if resolveAsString(path, left(n), info) == right(n).sval {
				return true
			}
		} else {
			if resolveAsInt(left(n), info) == resolveAsInt(right(n), info) {
				return true
			}
		}
	} else if n.ntype == T_NEQ {
		if left(n).ntype == T_NAME {
			if resolveAsString(path, left(n), info) != right(n).sval {
				return true
			}
		} else if left(n).ntype == T_PATH {
			if resolveAsString(path, left(n), info) != right(n).sval {
				return true
			}
		} else {
			if resolveAsInt(left(n), info) != resolveAsInt(right(n), info) {
				return true
			}
		}
	} else if n.ntype == T_OR {
		if evaluate(path, info, left(n)) || evaluate(path, info, right(n)) {
			return true
		}
	} else if n.ntype == T_AND {
		if evaluate(path, info, left(n)) && evaluate(path, info, right(n)) {
			return true
		}
	} else if n.ntype == T_ISFILE {
		return !info.IsDir()
	} else if n.ntype == T_ISDIR {
		return info.IsDir()
	} else if n.ntype == T_CONTAINS {
		return contains(left(n).ntype, right(n).sval, path, info)
	}
	return false
}

func resolveAsString(path string, n *tnode, info os.FileInfo) string {
	if n.ntype == T_NAME {
		return info.Name()
	} else if n.ntype == T_PATH {
		return path
	}
	return ""
}

func resolveAsInt(n *tnode, info os.FileInfo) int {
	if n.ntype == T_INTEGER {
		return n.ival
	} else if n.ntype == T_SIZE {
		return int(info.Size())
	} else if n.ntype == T_LITERAL {
		// try to parse string as timestamp/date
		t, err := time.Parse(TIMESTAMP_FORMAT, n.sval)
		if err == nil {
			return int(t.Unix())
		}
		d, err := time.Parse(DATE_FORMAT, n.sval)
		if err == nil {
			return int(d.Unix())
		}
		fmt.Println("Failed to parse string value: '" + n.sval +
			"'. If this is a date/time, be sure to follow one of the following two formats exactly: " +
			"\"" + TIMESTAMP_FORMAT + "\", \"" + DATE_FORMAT + "\"")
		os.Exit(1)
	} else if n.ntype == T_MODIFIED {
		_, zone := info.ModTime().Zone()
		modtime := int(info.ModTime().Unix()) + zone
		return modtime
	}
	return 0
}

func contains(ntype int, search string, path string, info os.FileInfo) bool {
	if ntype == T_NAME {
		return strings.Contains(info.Name(), search)
	} else if ntype == T_PATH  {
		return strings.Contains(path, search)
	} else if ntype == T_CONTENT {
		return fileContainsString(path, search)
	}
	return false
}

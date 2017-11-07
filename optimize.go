package main

import (
	"regexp"
)

const WEIGHT_NONE int = 0
const WEIGHT_INFO int = 1
const WEIGHT_FILE_STARTSWITH int = 2
const WEIGHT_FILE_ENDSWITH int = 3
const WEIGHT_FILE_MATCHES int = 4
const WEIGHT_FILE_CONTAINS int = 5

func shiftShortestPathLeft(n *tnode) int {
	if (n.ntype == T_CONTAINS || n.ntype == T_ICCONTAINS) && left(n).ntype == T_CONTENT {
		return WEIGHT_FILE_CONTAINS
	} else if (n.ntype == T_STARTSWITH || n.ntype == T_ICSTARTSWITH) && left(n).ntype == T_CONTENT {
		return WEIGHT_FILE_STARTSWITH
	} else if (n.ntype == T_ENDSWITH || n.ntype == T_ICENDSWITH) && left(n).ntype == T_CONTENT {
		return WEIGHT_FILE_ENDSWITH
	} else if (n.ntype == T_MATCHES && left(n).ntype == T_CONTENT) {
		return WEIGHT_FILE_MATCHES
	} else if n.ntype == T_AND || n.ntype == T_OR {
		// test each side and swap if necessary
		lval := shiftShortestPathLeft(left(n))
		rval := shiftShortestPathLeft(right(n))
		if lval > rval {
			// rhs has a shorter path than lhs, swap these
			tmp := left(n)
			n.children[0] = right(n)
			n.children[1] = tmp
			return rval
		} else {
			return lval
		}
	} else if n.ntype == T_NOT {
		return shiftShortestPathLeft(left(n))
	} else {
		return WEIGHT_NONE
	}
}

func compileRegexes(n *tnode) string {
	if n == nil {
		return ""
	}

	if n.ntype == T_MATCHES {
		var err error
		right(n).regval, err = regexp.Compile(right(n).sval)
		if err != nil {
			return right(n).sval
		}
	} else if n.ntype == T_AND || n.ntype == T_OR {
		str := compileRegexes(left(n))
		if str != "" {
			return str
		}

		str = compileRegexes(right(n))
		if str != "" {
			return str
		}
	}
	return ""
}

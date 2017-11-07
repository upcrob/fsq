package main

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
)

const (
	T_DEFAULT = iota
	T_PROGRAM
	T_OR
	T_AND
	T_NOT
	T_INTEGER
	T_CONTAINS
	T_ICCONTAINS
	T_CONTENT
	T_MODIFIED
	T_STARTSWITH
	T_ICSTARTSWITH
	T_ENDSWITH
	T_ICENDSWITH
	T_NAME
	T_PATH
	T_STATS
	T_ISFILE
	T_ISDIR
	T_SIZE
	T_FSIZE
	T_ALIST
	T_MATCHES
	T_IN
	T_LT
	T_LTE
	T_GT
	T_GTE
	T_EQ
	T_ICEQ
	T_NEQ
	T_ICNEQ
	T_STRING
	T_SLIST
)

type tnode struct {
	ntype    int
	ival     int
	sval     string
	regval   *regexp.Regexp
	children []*tnode
}

func left(n *tnode) *tnode {
	return n.children[0]
}

func right(n *tnode) *tnode {
	return n.children[1]
}

func addChild(n *tnode, c *tnode) {
	if n.children == nil {
		n.children = make([]*tnode, 0, 2)
	}
	n.children = append(n.children, c)
}

func printTree(root *tnode) {
	fmt.Println("PARSE TREE:")
	printTreeHelper(root, 0)
}

func printTreeHelper(n *tnode, depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(nodeString(n))
	for i := 0; i < len(n.children); i++ {
		printTreeHelper(n.children[i], depth+1)
	}
}

func collectFileSearchStrings(n *tnode) []SearchString {
	searchStrings := make([]SearchString, 0, 5)
	if n.ntype == T_CONTAINS ||
		n.ntype == T_STARTSWITH ||
		n.ntype == T_ENDSWITH {
		if n.children[0].ntype == T_CONTENT {
			searchStrings = append(searchStrings, SearchString{n.children[1].sval, true})
		}
	} else if n.ntype == T_ICCONTAINS ||
		n.ntype == T_ICSTARTSWITH ||
		n.ntype == T_ICENDSWITH {
		if n.children[0].ntype == T_CONTENT {
			searchStrings = append(searchStrings, SearchString{strings.ToLower(n.children[1].sval), false})
		}
	} else {
		for _, c := range n.children {
			rec := collectFileSearchStrings(c)
			for _, s := range rec {
				searchStrings = append(searchStrings, s)
			}
		}
	}
	return searchStrings
}

func nodeString(treeNode *tnode) string {
	switch treeNode.ntype {
	case T_PROGRAM:
		return "PROGRAM"
	case T_ALIST:
		return "ATTRIBUTE LIST"
	case T_OR:
		return "OR"
	case T_AND:
		return "AND"
	case T_INTEGER:
		return "INTEGER (" + strconv.Itoa(treeNode.ival) + ")"
	case T_STRING:
		return "STRING (\"" + treeNode.sval + "\")"
	case T_SLIST:
		return "STRING LIST"
	case T_CONTAINS:
		return "CONTAINS"
	case T_ICCONTAINS:
		return "ICCONTAINS"
	case T_CONTENT:
		return "CONTENT"
	case T_MODIFIED:
		return "MODIFIED"
	case T_STARTSWITH:
		return "STARTSWITH"
	case T_ICSTARTSWITH:
		return "ICSTARTSWITH"
	case T_ENDSWITH:
		return "ENDSWITH"
	case T_ICENDSWITH:
		return "ICENDSWITH"
	case T_NAME:
		return "NAME"
	case T_PATH:
		return "PATH"
	case T_SIZE:
		return "SIZE"
	case T_FSIZE:
		return "FSIZE"
	case T_STATS:
		return "STATS"
	case T_MATCHES:
		return "MATCHES"
	case T_IN:
		return "IN"
	case T_LT:
		return "LT"
	case T_LTE:
		return "LTE"
	case T_GT:
		return "GT"
	case T_GTE:
		return "GTE"
	case T_EQ:
		return "EQ"
	case T_ICEQ:
		return "ICEQ"
	case T_NEQ:
		return "NEQ"
	case T_ICNEQ:
		return "ICNEQ"
	case T_ISFILE:
		return "ISFILE"
	case T_ISDIR:
		return "ISDIR"
	default:
		return "UNKNOWN NODE"
	}
}

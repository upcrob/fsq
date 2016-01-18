package main

// moves "contents contains" conditions toward the right
// side of the parse tree in order to avoid unnecessary
// file reading
func optimizeContentsContains(n *tnode) bool {
	if n == nil {
		return false
	}

	if n.ntype == T_OR {
		l := optimizeContentsContains(left(n))
		r := optimizeContentsContains(right(n))
		if l && !r {
			tmp := left(n)
			n.children[0] = right(n)
			n.children[1] = tmp
		}
		return l || r
	} else if n.ntype == T_AND {
		l := optimizeContentsContains(left(n))
		r := optimizeContentsContains(right(n))
		if l && !r {
			tmp := left(n)
			n.children[0] = right(n)
			n.children[1] = tmp
		}
		return l || r
	} else if n.ntype == T_CONTAINS && left(n).ntype == T_CONTENTS {
		return true
	}
	return false
}

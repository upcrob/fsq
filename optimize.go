package main

func isContentContainsExpression(n *tnode) bool {
	return n.ntype == T_CONTAINS && left(n).ntype == T_CONTENT
}

func isContentStartswithExpression(n *tnode) bool {
	return n.ntype == T_STARTSWITH && left(n).ntype == T_CONTENT
}

func isContentEndswithExpression(n *tnode) bool {
	return n.ntype == T_ENDSWITH && left(n).ntype == T_CONTENT
}

/*
 * Moves matching expressions toward the right-hand side
 * of the parse tree in order to avoid unnecessary
 * file reading.
 */
func shiftExpressionRight(n *tnode, isMatch func(n *tnode) bool) bool {
	if n == nil {
		return false
	}

	if n.ntype == T_OR || n.ntype == T_AND {
		l := shiftExpressionRight(left(n), isMatch)
		r := shiftExpressionRight(right(n), isMatch)
		if l && !r {
			tmp := left(n)
			n.children[0] = right(n)
			n.children[1] = tmp
		}
		return l || r
	}
	return isMatch(n)
}

package main

import (
	"strconv"
	"strings"
)

type Lexer struct {
	expr string
}

var keywordMappings map[int]string = map[int]string {
	NAME: "name",
	SIZE: "size",
	ISFILE: "isfile",
	ISDIR: "isdir",
	PATH: "path",
	WHERE: "where",
	IN: "in",
	CONTAINS: "contains",
	CONTENTS: "contents",
	AND: "and",
	OR: "or",
	STARTSWITH: "startswith",
	ENDSWITH: "endswith",
}


func (lexer *Lexer) Lex(lval *yySymType) int {
	// trim leading whitespace
	lexer.expr = lexer.expr[getWhitespaceCount(lexer.expr):]

	if lexer.expr == "" {
		return EOF
	} else if lexer.expr[0] == ',' {
		lexer.expr = lexer.expr[1:]
		return COMMA
	} else if lexer.expr[0] == '(' {
		lexer.expr = lexer.expr[1:]
		return OPAREN
	} else if lexer.expr[0] == ')' {
		lexer.expr = lexer.expr[1:]
		return CPAREN
	} else if sval := getInteger(lexer.expr); sval != "" {
		lval.ival, _ = strconv.Atoi(sval)
		lexer.expr = lexer.expr[len(sval):]
		return INTEGER
	} else if op, olen := getOperator(lexer.expr); op != -1 {
		lexer.expr = lexer.expr[olen:]
		return op
	} else if keyword, klen := getKeyword(lexer.expr); keyword != -1 {
		lexer.expr = lexer.expr[klen:]
		return keyword
	} else if literal, llen := getLiteral(lexer.expr); literal != "" {
		lexer.expr = lexer.expr[llen:]
		lval.sval = literal
		return LITERAL
	}

	lexer.expr = lexer.expr[1:]
	return UNKNOWN
}

func (lexer *Lexer) Error(e string) {

}

func tokenString(symbolId int) string {
	if str := keywordMappings[symbolId]; str != "" {
		return strings.ToUpper(str)
	}

	switch symbolId {
	case EOF:
		return "EOF"
	case COMMA:
		return "COMMA"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case LITERAL:
		return "LITERAL"
	case INTEGER:
		return "INTEGER"
	case STARTSWITH:
		return "STARTSWITH"
	case ENDSWITH:
		return "ENDSWITH"
	case ISDIR:
		return "ISDIR"
	case ISFILE:
		return "ISFILE"
	case OPAREN:
		return "OPAREN"
	case CPAREN:
		return "CPAREN"
	case PATH:
		return "PATH"
	default:
		return "UNKNOWN"
	}
}

func trim(s *string) {
	q := *s
	*s = q[1:]
}

func getWhitespaceCount(expr string) int {
	i := 0
	for i < len(expr) && (expr[i] == ' ' || expr[i] == '\n' || expr[i] == '\t') {
		i++
	}
	return i
}

func getLiteral(expr string) (string, int) {
	i := 0
	if expr[i] == '\'' {
		i++

		for i < len(expr) && expr[i] != '\'' {
			i++
		}

		if i < len(expr) && expr[i] == '\'' {
			return expr[1:i], i+1
		}
	}
	return "", 0
}

func getKeyword(expr string) (int, int) {
	ident := getIdent(expr)

	if ident != "" {
		for sym, str := range keywordMappings {
			if ident == str {
				return sym, len(str)
			}
		}
	}
	return -1, 0
}

func getInteger(expr string) string {
	i := 0
	for i < len(expr) && isNumeric(expr[i]) {
		i++
	}

	if i > 0 {
		return expr[:i]
	}
	return ""
}

func getIdent(expr string) string {
	i := 0
	for i < len(expr) && isAlpha(expr[i]) {
		i++
	}

	if i > 0 {
		return expr[:i]
	}
	return ""
}

func getOperator(expr string) (int, int) {
	i := 0
	for i < len(expr) && (expr[i] == '>' || expr[i] == '<' || expr[i] == '=' || expr[i] == '!') {
		i++
	}

	if i > 0 {
		op := expr[:i]
		l := len(op)

		switch op {
		case ">":
			return GT, l
		case "<":
			return LT, l
		case ">=":
			return GTE, l
		case "<=":
			return LTE, l
		case "=":
			return EQ, l
		case "!=":
			return NEQ, l
		}
		return UNKNOWN, l
	}
	return -1, 0
}

func isAlpha(c byte) bool {
	if c >= 65 && c <= 90 || c >= 97 && c <= 122 {
		return true
	} else {
		return false
	}
}

func isNumeric(c byte) bool {
	if c >= 48 && c <= 57 {
		return true
	} else {
		return false
	}
}
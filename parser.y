%{
	package main

	var programRoot *tnode
%}

%union {
	ival int
	sval string
	tval *tnode
}

%token UNKNOWN
%token EOF
%token COMMA
%token OPAREN
%token CPAREN
%token <sval> LITERAL
%token <sval> NAME
%token <ival> SIZE
%token <sval> PATH
%token ISDIR
%token ISFILE
%token WHERE
%token IN
%token AND
%token OR
%token CONTAINS
%token CONTENTS
%token STARTSWITH
%token ENDSWITH
%token LT
%token LTE
%token GT
%token GTE
%token EQ
%token NEQ
%token <ival> INTEGER

%type <tval> attribute attribute_list location or_expr and_expr logic_expr value program

// ========================================================
// BEGIN GRAMMAR
// ========================================================

%%

program:
	attribute_list location LITERAL WHERE or_expr EOF {
		programRoot = new(tnode)
		programRoot.ntype = T_PROGRAM
		addChild(programRoot, $1)
		addChild(programRoot, $2)
		n := new(tnode)
		n.ntype = T_LITERAL
		n.sval = $3
		addChild(programRoot, n)
		addChild(programRoot, $5)
	}
	;

value:
	LITERAL {
		$$ = new(tnode)
		$$.ntype = T_LITERAL
		$$.sval = $1
	}
	| INTEGER {
		$$ = new(tnode)
		$$.ntype = T_INTEGER
		$$.ival = $1
	}
	;

or_expr:
	or_expr OR and_expr {
		$$ = new(tnode)
		$$.ntype = T_OR
		addChild($$, $1)
		addChild($$, $3)
	}
	| and_expr {
		$$ = $1
	}
	;

and_expr:
	and_expr AND logic_expr {
		$$ = new(tnode)
		$$.ntype = T_AND
		addChild($$, $1)
		addChild($$, $3)
	}
	| logic_expr {
		$$ = $1
	}
	;

logic_expr:
	ISDIR {
		$$ = new(tnode)
		$$.ntype = T_ISDIR
	}
	| ISFILE {
		$$ = new(tnode)
		$$.ntype = T_ISFILE
	}
	| attribute LT value {
		$$ = new(tnode)
		$$.ntype = T_LT
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute LTE value {
		$$ = new(tnode)
		$$.ntype = T_LTE
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute GT value {
		$$ = new(tnode)
		$$.ntype = T_GT
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute GTE value {
		$$ = new(tnode)
		$$.ntype = T_GTE
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute EQ value {
		$$ = new(tnode)
		$$.ntype = T_EQ
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute NEQ value {
		$$ = new(tnode)
		$$.ntype = T_NEQ
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute CONTAINS value {
		$$ = new(tnode)
		$$.ntype = T_CONTAINS
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute STARTSWITH value {
		$$ = new(tnode)
		$$.ntype = T_STARTSWITH
		addChild($$, $1)
		addChild($$, $3)
	}
	| attribute ENDSWITH value {
		$$ = new(tnode)
		$$.ntype = T_ENDSWITH
		addChild($$, $1)
		addChild($$, $3)
	}
	| OPAREN or_expr CPAREN {
		$$ = $2
	}
	;

location:
	IN {
		$$ = new(tnode)
		$$.ntype = T_IN
	}
	;

attribute_list:
	attribute {
		$$ = new(tnode)
		$$.ntype = T_ALIST
		addChild($$, $1)
	}
	| attribute COMMA attribute_list {
		$$ = $3
		addChild($3, $1)
	}
	;

attribute:
	NAME {
		$$ = new(tnode)
		$$.ntype = T_NAME
	}
	| SIZE {
		$$ = new(tnode)
		$$.ntype = T_SIZE
	}
	| PATH {
		$$ = new(tnode)
		$$.ntype = T_PATH
	}
	| CONTENTS {
		$$ = new(tnode)
		$$.ntype = T_CONTENTS
	}
	;
%%

%{  /* prologue section of code defining C/C++ code to directly insert in parser */
    #include "astnode.hpp"
    BlockNode* program_block;

    extern int yylex();  /*  use flex function to get tokens */
    void yyerror(const char* s) { printf("error: %s", s); }
%}

/* different ways to access token data */
%union {
    Node* node;
    BlockNode* block;
    StmtNode* stmt;
    ExprNode* expr;
    IdentifierNode* ident;
    VarDeclNode* var_decl;

    std::vector<VarDeclNode*>* var_decl_list;
    std::vector<ExprNode*>* expr_list;

    std::string* string;
    int token;
}

/* define terminal symbols, matching with those defined in lexer */
%token <string> IDENTIFIER INTEGER DOUBLE
%token <token> PLUS MINUS MUL DIV
%token <token> LPAREN RPAREN LBRACE RBRACE DOT COMMA SEMICOLON
%token <token> LT GT ASSIGN EQUAL

/* define non-terminal symbols */
%type <ident> ident
%type <expr> numeric expr
%type <var_decl_list> func_decl_args
%type <expr_list> call_args
%type <block> program stmts block
%type <stmt> stmt var_decl func_decl

/* define operator precedence */
%right ASSIGN
%left EQUAL
%left LT GT
%left PLUS MINUS
%left MUL DIV

/* define start symbol of grammar */
%start program

%%

/* define CFG and production rules */

program 
    : stmts { program_block = $1; }
    ;

stmts 
    : stmt { $$ = new BlockNode(); $$->stmts.push_back($<stmt>1); }
    | stmts stmt { $1->stmts.push_back($<stmt>2); }
    ;

stmt
    : var_decl | func_decl
    | expr SEMICOLON { $$ = new ExprStmtNode(*$1); }
    ;

block
    : LBRACE stmts RBRACE { $$ = $2; }
    | LBRACE RBRACE { $$ = new BlockNode(); }
    ;

ident
    : IDENTIFIER { $$ = new IdentifierNode(*$1); delete $1; }
    ;

var_decl
    : ident ident { $$ = new VarDeclNode(*$1, *$2); }
    | ident ident ASSIGN expr { $$ = new VarDeclNode(*$1, *$2, $4); }
    ;

func_decl_args
    : /* blank */ {$$ = new VarDeclList(); }
    | var_decl { $$ = new VarDeclList(); $$->push_back($<var_decl>1); }
    | func_decl_args COMMA var_decl { $1->push_back($<var_decl>3); }
    ;

func_decl
    : ident ident LPAREN func_decl_args RPAREN block {
        $$ = new FuncDeclNode(*$1, *$2, *$4, *$6); delete $4; 
    }
    ;

numeric
    : INTEGER { $$ = new IntegerNode(atol($1->c_str())); delete $1; }
    | DOUBLE { $$ = new DoubleNode(atof($1->c_str())); delete $1; }
    ;

expr
    : ident ASSIGN expr {$$ = new AssignmentNode(*$<ident>1, *$3); }
    | ident LPAREN call_args RPAREN { $$ = new MethodCallNode(*$1, *$3); delete $3; }
    | ident { $<ident>$ = $1; }
    | numeric
        | expr PLUS expr { $$ = new BinaryOpNode(*$1, *$3, PLUS); }
        | expr MINUS expr { $$ = new BinaryOpNode(*$1, *$3, MINUS); }
        | expr MUL expr { $$ = new BinaryOpNode(*$1, *$3, MUL); }
        | expr DIV expr { $$ = new BinaryOpNode(*$1, *$3, DIV); }
        | expr LT expr { $$ = new BinaryOpNode(*$1, *$3, LT); }
        | expr GT expr { $$ = new BinaryOpNode(*$1, *$3, GT); }
    | expr EQUAL expr { $$ = new BinaryOpNode(*$1, *$3, EQUAL); }
    | LPAREN expr RPAREN { $$ = $2; }
    ;

call_args
    : /* blank */ { $$ = new ExprList(); }
    | expr { $$ = new ExprList(); $$->push_back($1); }
    | call_args COMMA expr { $1->push_back($3); }
    ;

%%

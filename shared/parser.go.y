%{
package trompe
%}

%union{
    tok Token
    node *Node
    nodelist []*Node
    word *Word
    wordlist []*Word
}

%token<word> WILDCARD LIDENT UIDENT CHAR INT FLOAT STRING REGEXP SOME NONE
%token<tok> ABSTRACT AND AS ASSERT BEGIN CONSTRAINT DO DONE DOWNTO ELSE END EXCEPTION EXTERNAL FALSE FOR FUN FUNCTION GOTO IF IMPORT IN LET MATCH MOD MODULE MUTABLE NOT OF OPEN OR OVERRIDE PARTIAL REC RETURN SIG STRUCT THEN TO TRAIT TRUE TRY TYPE USE VAL WHEN WHILE WITH WITHOUT
%token<tok> ADD ADD_DOT SUB SUB_DOT MUL MUL_DOT DIV DIV_DOT REM POW EQ NE LE GE LT GT LPAREN RPAREN LBRACE RBRACE LBRACK RBRACK DCOLON SEMI SEMI2 COLON COMMA DOT DOT2 DOT3  DOT_LPAREN QUOTE Q EP AMP TILDA OP
%token<tok> DOL PIPE BAND BOR BXOR LSHIFT RSHIFT LARROW RARROW
%token<tok> INDENT DEDENT EOF

%type<node> program def let pattern exp binexp param arg use excdef
%type<node> simple_exp valuepath valuename constr constant tuple list array seqexp match ptnmatch_head
%type<node> typdef typeq typexp raw_typexp simple_typexp label modtyp multimatch
%type<node> modpath typconstr
%type<nodelist> moditems deflist letbind args tuple_comps trait_params patterncomp patternlist eltlist eltarray paramlist ptnmatch ptnmatch_tail seqexplist typexplist_comma typexplist_ast
%type<word> modname constrname modtypname typconstrname fieldname prefix
%type<wordlist> vallist

%right SEMI2
%nonassoc prec_let prec_fun prec_try prec_raw_typexp
%right SEMI
%right prec_if
%left DEDENT
%right ELSE
%nonassoc AS
%right RARROW
%right prec_ptnmatch
%left PIPE
%nonassoc COMMA
%right DOL
%left LOR
%left LAND
%left EQ NE LG LT GT LE GE
%right CONCAT
%right COLON2
%right prec_constr
%left ADD SUB ADD_DOT SUB_DOT
%nonassoc prec_typexplist_ast
%left MUL MUL_DOT DIV DIV_DOT MOD BAND BOR BXOR
%right LSHIFT RSHIFT
%right prec_unary_minus prec_unary_minus_dot
%left prec_app prec_assert
%nonassoc DOT DOT_LPAREN
%nonassoc prec_prefix

%right UIDENT
%nonassoc prec_simple_typexp
%nonassoc LIDENT
%nonassoc COLON
%left LPAREN LBRACE LBRACK

%nonassoc RPAREN

%nonassoc prec_simple_exp
%right prec_kw

%nonassoc INT FLOAT BOOL CHAR STRING TRUE FALSE BEGIN

%%

program
    : moditems
    {
        $$ = newNode($1[0].Loc, &ProgramNode{Items:$1})
        if p, ok := yylex.(*Parser); ok {
            p.unit = $$
        }
    }

moditems
    : deflist { $$ = $1 }

deflist
    : def { $$ = NewNodeList($1) }
    | deflist def { $$ = append($1, $2) }

def
    /* interface only */
    : VAL valuename COLON typexp
    { $$ = newNode($1.Loc, &ValueDefNode{Name:$2, TypeExp:$4}) }

    /* interface and implementation */
    | EXTERNAL valuename COLON typexp EQ STRING
    { $$ = newNode($1.Loc.Union($6.Loc), &ExtNode{Name:$2.Desc.(*IdentNode).Name, TypeExp:$4, Prim:$6.Value}) }
    | typdef { $$ = $1 }
    | excdef { $$ = $1 }
    | MODULE modname COLON modtyp
    { $$ = newNode($1.Loc, &ModuleDefNode{Name:$2, TypeExp:$4}) }
    | TRAIT modname COLON modtyp
    { $$ = newNode($1.Loc, &TraitDefNode{Name:$2, TypeExp:$4}) }
    | IMPORT modpath { $$ = newNode($1.Loc, &ImportNode{Path:$2}) }
    | IMPORT modpath AS UIDENT
    { $$ = newNode($1.Loc, &ImportNode{Path:$2, Alias:$4}) }
    | OPEN IMPORT modpath { $$ = newNode($1.Loc, &ImportNode{Open:true, Path:$3}) }
    | OPEN IMPORT modpath AS UIDENT
    { $$ = newNode($1.Loc, &ImportNode{Open:true, Path:$3, Alias:$5}) }

    /* implementation only */
    | LET letbind
    { $$ = newNode($1.Loc, &LetNode{Public:true, Rec:false, Bindings:$2}) }
    | LET REC letbind
    {
        for _, bind := range $3 {
            if desc,ok := bind.Desc.(*LetBindingNode);ok{
                desc.Rec = true
            } else if desc,ok := bind.Desc.(*BlockNode);ok{
                desc.Rec = true
            } else {
                panic("error")
            }
        }
        $$ = newNode($1.Loc, &LetNode{Public:true, Rec:true, Bindings:$3})
    }
    | SEMI2 seqexp { $$ = $2 }
    | use { $$ = $1 }

constrdecl
    : constrname {}
    | constrname OF simple_typexp {}
    | constrname OF simple_typexp typexplist_ast {}

modtyp
    : modtyppath {}
    | SIG deflist END {}
    | modtyp WITH modconstraintlist_and {}
    | LPAREN modtyp RPAREN {}

modtyppath
    : modtypname {}
    | modpath DOT modtypname {}

modtypname
    : UIDENT { $$ = $1 }

modconstraintlist_and
    : modconstraint {}
    | modconstraintlist_and AND modconstraint {}

modconstraint
    : TYPE typconstr typeq {}
    | TYPE typparamlist typconstr typeq {}
    | MODULE modpath EQ modpath {}

typconstr
    : typconstrname { $$ = newNode($1.Loc, &TypeConstrNode{Name:$1}) }
    | modpath DOT typconstrname
    { $$ = newNode($1.Loc, &TypeConstrNode{Path:$1, Name:$3}) }

typconstrname
    : LIDENT { $$ = $1 }

excdef
    : EXCEPTION constrdecl {}
    | EXCEPTION constrname EQ constr {}

constr
    : constrname { $$ = newNode($1.Loc, &ConstrNode{Name:$1}) }
    | modpath DOT constrname
    { $$ = newNode($1.Loc, &ConstrNode{Path:$1, Name:$3}) }

constrname
    : UIDENT { $$ = $1 }

letbind
    : let { $$ = NewNodeList($1) }
    | letbind AND let { $$ = append($1, $3) }

let
    : pattern EQ seqexp
    { $$ = newNode($1.unionLoc($3), &LetBindingNode{Ptn:$1, Body:$3}) }
    | pattern EQ INDENT seqexp DEDENT
    { $$ = newNode($1.unionLoc($4), &LetBindingNode{Ptn:$1, Body:$4}) }
    | valuename paramlist EQ seqexp
    { $$ = newNode($1.unionLoc($4), &BlockNode{Name:$1, Params:$2, Body: $4}) }
    | valuename paramlist EQ INDENT seqexp DEDENT
    { $$ = newNode($1.unionLoc($5), &BlockNode{Name:$1, Params:$2, Body: $5}) }

pattern
    : valuename
    { $$ = newNode($1.Loc, &PtnIdentNode{Name:$1.Desc.(*IdentNode).Name}) }
    | WILDCARD
    { $$ = newNode($1.Loc, &WildcardNode{}) }
    | constant
    { $$ = newNode($1.Loc, &PtnConstNode{Value:$1}) }
    | constr pattern %prec prec_constr
    { $$ = newNode($1.Loc, &PtnConstrAppNode{Constr:$1, Ptn:$2}) }
    | SOME pattern %prec prec_constr
    { $$ = newNode($1.Loc, &PtnSomeNode{Ptn:$2}) }
    | pattern AS valuename
    { $$ = newNode($1.unionLoc($3), &PtnVarNode{Ptn:$1, Name:$3.Desc.(*IdentNode).Name}) }
    | LPAREN pattern COLON typexp RPAREN
    { $$ = newNode($1.Loc.Union($5.Loc), &PtnTypeNode{Ptn:$2, TypeExp:$4}) }
    | pattern PIPE pattern
    { $$ = newNode($1.unionLoc($3), &SeqPtnNode{Left:$1, Right:$3}) }
    | LBRACK RBRACK
    { $$ = newNode($1.Loc, &PtnListNode{}) }
    | LBRACK patternlist RBRACK
    { $$ = newNode($2[0].Loc, &PtnListNode{Elts:$2}) }
    | LBRACK patternlist SEMI RBRACK
    { $$ = newNode($2[0].Loc, &PtnListNode{Elts:$2}) }
    | pattern COLON2 pattern
    { $$ = newNode($1.unionLoc($3), &PtnListConsNode{Head:$1, Tail:$3}) }
    | LPAREN patterncomp RPAREN
    { $$ = newNode($2[0].Loc, &PtnTupleNode{Comps:$2}) }

patternlist
    : pattern { $$ = NewNodeList($1) }
    | patternlist SEMI pattern { $$ = append($1, $3) }

patterncomp
    : pattern { $$ = NewNodeList($1) }
    | patterncomp COMMA pattern { $$ = append($1, $3) }

valuename
    : LIDENT { $$ = newNode($1.Loc, &IdentNode{Name:$1.Value}) }
    | LPAREN ADD RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"+"}) }
    | LPAREN SUB RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"-"}) }
    | LPAREN MUL RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"*"}) }
    | LPAREN DIV RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"/"}) }
    | LPAREN MOD RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"mod"}) }
    | LPAREN REM RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"%"}) }
    | LPAREN POW RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"**"}) }
    | LPAREN ADD_DOT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"+."}) }
    | LPAREN SUB_DOT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"-."}) }
    | LPAREN MUL_DOT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"*."}) }
    | LPAREN DIV_DOT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"/."}) }
    | LPAREN EQ RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"="}) }
    | LPAREN NE RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"<>"}) }
    | LPAREN LT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"<"}) }
    | LPAREN LE RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"<="}) }
    | LPAREN GT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:">"}) }
    | LPAREN GE RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:">="}) }
    | LPAREN CONCAT RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"^"}) }
    | LPAREN PIPE RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"|"}) }
    | LPAREN DOL RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"$"}) }
    | LPAREN AMP RPAREN { $$ = newNode($1.Loc, &IdentNode{Name:"&"}) }

paramlist
    : param { $$ = NewNodeList($1) }
    | paramlist param { $$ = append($1, $2) }

param
    : pattern { $$ = $1 }
    | LIDENT COLON pattern
    { $$ = newNode($1.Loc, &LabelParamNode{Name:$1, Ptn:$3}) }
    | COLON LIDENT
    {
        ptn := newNode($1.Loc, &PtnIdentNode{Name:$2.Value})
        $$ = newNode($1.Loc, &LabelParamNode{Name:$2, Ptn:ptn})
    }

/* TODO
    | KEYWORD {}
    | LPAREN KEYWORD RPAREN {}
    | LPAREN KEYWORD COLON typexp RPAREN {}
    | KEYWORD COLON pattern {}
    | Q KEYWORD {}
    | Q RPAREN KEYWORD RPAREN {}    
    | Q RPAREN KEYWORD COLON typexp RPAREN {}    
    | Q RPAREN KEYWORD COLON typexp EQ exp RPAREN {}    
    | Q RPAREN KEYWORD EQ exp RPAREN {}    
    | Q KEYWORD COLON pattern {}
    | Q KEYWORD COLON RPAREN pattern RPAREN {}
    | Q KEYWORD COLON RPAREN pattern COLON typexp RPAREN {}
    | Q KEYWORD COLON RPAREN pattern COLON typexp EQ exp RPAREN {}
    | Q KEYWORD COLON RPAREN pattern EQ exp RPAREN {}
*/

seqexp
    : seqexplist { $$ = newNode($1[0].Loc, &SeqExpNode{Exps:$1}) }

seqexplist
    : exp { $$ = NewNodeList($1) }
    | seqexplist SEMI exp { $$ = append($1, $3) }

exp
    : simple_exp { $$ = $1 }
    | binexp { $$ = $1 }
    | constr simple_exp
    { $$ = newNode($1.Loc, &ConstrAppNode{Constr:$1, Exp:$2}) }
    | simple_exp COLON2 simple_exp
    { $$ = newNode($1.unionLoc($3), &ListConsNode{Head:$1, Tail:$3}) }
    | prefix exp %prec prec_prefix
    { $$ = newNode($1.Loc, &PrefixExpNode{Prefix:$1, Exp:$2}) }
    | simple_exp args %prec prec_app
    { $$ = newNode($1.Loc, &AppNode{Exp:$1, Args:$2}) }
    | LET letbind IN exp %prec prec_let
    { $$ = newNode($1.Loc, &LetNode{Public:false, Rec:false, Bindings:$2, Body:$4}) }
    | LET REC letbind IN exp %prec prec_let
    { $$ = newNode($1.Loc, &LetNode{Public:false, Rec:true, Bindings:$3, Body:$5}) }
    | IF exp THEN exp %prec prec_if
    { $$ = newNode($1.Loc.Union($4.Loc), &IfNode{Cond:$2, True:$4}) }
    | IF exp THEN exp ELSE exp %prec prec_if
    { $$ = newNode($1.Loc.Union($6.Loc), &IfNode{Cond:$2, True:$4, False:$6}) }
    | IF exp THEN INDENT seqexp DEDENT ELSE exp
    { $$ = newNode($1.Loc.Union($8.Loc), &IfNode{Cond:$2, True:$5, False:$8}) }
    | MATCH exp WITH ptnmatch
    { $$ = newNode($1.Loc.Union($3.Loc), &CaseNode{Exp:$2, Match:$4}) }
    | TRY exp WITH ptnmatch
    { $$ = newNode($1.Loc.Union($3.Loc), &TryNode{Exp:$2, Match:$4}) }
    | FUNCTION ptnmatch
    { $$ = newNode($1.Loc, &FunctionNode{Match:$2}) }
    | FUN multimatch %prec prec_fun
    { $$ = newNode($1.Loc, &FunNode{MultiMatch:$2}) }
    | SOME simple_exp { $$ = newNode($1.Loc, &OptionNode{Value:$2}) }
    | simple_exp DOT_LPAREN exp RPAREN LARROW simple_exp
    { $$ = newNode($1.Loc, &ArrayAccessNode{Array:$1, Index:$3, Set: $6}) }

prefix
    : EP { $$ = &Word{Loc:$1.Loc, Value:"!"} }
    | Q { $$ = &Word{Loc:$1.Loc, Value:"?"} }
    | TILDA { $$ = &Word{Loc:$1.Loc, Value:"~"} }
    | SUB { $$ = &Word{Loc:$1.Loc, Value:"-"} }
    | SUB_DOT { $$ = &Word{Loc:$1.Loc, Value:"-."} }

ptnmatch
    : ptnmatch_head %prec prec_ptnmatch
    { $$ = NewNodeList($1) }
    | ptnmatch_head ptnmatch_tail %prec prec_ptnmatch
    {
        $$ = NewNodeList($1)
        for _, m := range $2 {
            $$ = append($$, m)
        }
    }

ptnmatch_tail
    : match { $$ = NewNodeList($1) }
    | ptnmatch_tail match { $$ = append($1, $2) }

ptnmatch_head
    : pattern RARROW exp
    { $$ = newNode($1.Loc.Union($3.Loc), &MatchNode{Ptn:$1, Body:$3}) }
    | pattern RARROW INDENT exp DEDENT
    { $$ = newNode($1.Loc.Union($4.Loc), &MatchNode{Ptn:$1, Body:$4}) }
    | pattern WHEN exp RARROW exp
    { $$ = newNode($1.Loc.Union($5.Loc), &MatchNode{Ptn:$1, Cond:$3, Body:$5}) }
    | pattern WHEN exp RARROW INDENT exp DEDENT
    { $$ = newNode($1.Loc.Union($6.Loc), &MatchNode{Ptn:$1, Cond:$3, Body:$6}) }
    | PIPE pattern RARROW exp
    { $$ = newNode($1.Loc.Union($4.Loc), &MatchNode{Ptn:$2, Body:$4}) }
    | PIPE pattern RARROW INDENT exp DEDENT
    { $$ = newNode($1.Loc.Union($5.Loc), &MatchNode{Ptn:$2, Body:$5}) }
    | PIPE pattern WHEN exp RARROW exp
    { $$ = newNode($1.Loc.Union($6.Loc), &MatchNode{Ptn:$2, Cond:$4, Body:$6}) }
    | PIPE pattern WHEN exp RARROW INDENT exp DEDENT
    { $$ = newNode($1.Loc.Union($7.Loc), &MatchNode{Ptn:$2, Cond:$4, Body:$7}) }

match
    : PIPE pattern RARROW exp
    { $$ = newNode($1.Loc.Union($4.Loc), &MatchNode{Ptn:$2, Body:$4}) }
    | PIPE pattern RARROW INDENT exp DEDENT
    { $$ = newNode($1.Loc.Union($5.Loc), &MatchNode{Ptn:$2, Body:$5}) }
    | PIPE pattern WHEN exp RARROW exp
    { $$ = newNode($1.Loc.Union($6.Loc), &MatchNode{Ptn:$2, Cond:$4, Body:$6}) }
    | PIPE pattern WHEN exp RARROW INDENT exp DEDENT
    { $$ = newNode($1.Loc.Union($7.Loc), &MatchNode{Ptn:$2, Cond:$4, Body:$7}) }

multimatch 
    : paramlist RARROW exp
    { $$ = newNode($1[0].Loc, &MultiMatchNode{Params:$1, Body:$3}) }
    | paramlist RARROW INDENT exp DEDENT
    { $$ = newNode($1[0].Loc, &MultiMatchNode{Params:$1, Body:$4}) }
    | paramlist WHEN exp RARROW exp
    { $$ = newNode($1[0].Loc, &MultiMatchNode{Params:$1, Cond:$3, Body:$5}) }
    | paramlist WHEN exp RARROW INDENT exp DEDENT
    { $$ = newNode($1[0].Loc, &MultiMatchNode{Params:$1, Cond:$3, Body:$6}) }

simple_exp
    : valuepath { $$ = $1 }
    | constant { $$ = $1 }
    | list { $$ = $1 }
    | tuple { $$ = $1 }
    | array { $$ = $1 }
    | LPAREN exp RPAREN { $$ = $2 }
    | LPAREN INDENT exp RPAREN { $$ = $3 }
    | LPAREN INDENT exp DEDENT RPAREN { $$ = $3 }
    | BEGIN exp END { $$ = $2 }
    | BEGIN INDENT exp DEDENT { $$ = $3 }
    | LPAREN exp COLON typexp RPAREN
    { $$ = newNode($1.Loc.Union($5.Loc), &TypeSpecifiedExpNode{Exp:$2, TypeExp:$4}) }
    | IF exp THEN INDENT seqexp DEDENT
    { $$ = newNode($1.Loc.Union($5.Loc), &IfNode{Cond:$2, True:$5}) }
    | IF exp THEN INDENT seqexp DEDENT ELSE INDENT seqexp DEDENT
    { $$ = newNode($1.Loc.Union($9.Loc), &IfNode{Cond:$2, True:$5, False:$9}) }
    | WHILE exp DO seqexp DONE
    { $$ = newNode($1.Loc.Union($5.Loc), &WhileNode{Cond:$2, Body:$4}) }
    | WHILE exp DO INDENT seqexp DEDENT
    { $$ = newNode($1.Loc.Union($5.Loc), &WhileNode{Cond:$2, Body:$5}) }
    | FOR valuename EQ exp TO exp DO seqexp DONE
    { $$ = newNode($1.Loc.Union($9.Loc), &ForNode{Name:$2, Init:$4, Dir:ForDirTo, Limit:$6, Body:$8}) }
    | FOR valuename EQ exp TO exp DO INDENT seqexp DEDENT
    { $$ = newNode($1.Loc.Union($9.Loc), &ForNode{Name:$2, Init:$4, Dir:ForDirTo, Limit:$6, Body:$9}) }
    | FOR valuename EQ exp DOWNTO exp DO seqexp DONE
    { $$ = newNode($1.Loc.Union($9.Loc), &ForNode{Name:$2, Init:$4, Dir:ForDirDownTo, Limit:$6, Body:$8}) }
    | FOR valuename EQ exp DOWNTO exp DO INDENT seqexp DEDENT
    { $$ = newNode($1.Loc.Union($9.Loc), &ForNode{Name:$2, Init:$4, Dir:ForDirDownTo, Limit:$6, Body:$9}) }
    | MATCH exp WITH INDENT ptnmatch DEDENT
    { $$ = newNode($1.Loc.Union($5[0].Loc), &CaseNode{Exp:$2, Match:$5}) }
    | TRY exp WITH INDENT ptnmatch DEDENT
    { $$ = newNode($1.Loc.Union($5[0].Loc), &TryNode{Exp:$2, Match:$5}) }
    | simple_exp DOT_LPAREN exp RPAREN
    { $$ = newNode($1.Loc, &ArrayAccessNode{Array:$1, Index:$3}) }

valuepath
    : valuename { $$ = $1 }
    | modpath DOT valuename
    { $$ = newNode($1.Loc, &ValuePathNode{Path:$1, Name:$3}) }

modpath
    : modname { $$ = newNode($1.Loc, &ModulePathNode{Path:NewWordList($1)}) }
    | modpath DOT modname
    {
        desc := $1.Desc.(*ModulePathNode)
        desc.Path = append(desc.Path, $3)
        $$ = $1
    }

modname
    : UIDENT { $$ = $1 }

constant
    : LPAREN RPAREN { $$ = newNode($1.Loc.Union($2.Loc), &UnitNode{}) }
    | TRUE { $$ = newNode($1.Loc, &BoolNode{true}) }
    | FALSE { $$ = newNode($1.Loc, &BoolNode{false}) }
    | INT { $$ = newNode($1.Loc, &IntNode{$1.Value}) }
    | FLOAT { $$ = newNode($1.Loc, &FloatNode{$1.Value}) }
    | STRING { $$ = newNode($1.Loc, &StringNode{$1.Value}) }
    | REGEXP { $$ = newNode($1.Loc, &RegexpNode{$1.Value}) }
    | NONE { $$ = newNode($1.Loc, &OptionNode{}) }

list
    : LBRACK RBRACK { $$ = newNode($1.Loc, &ListNode{}) } 
    | LBRACK eltlist RBRACK { $$ = newNode($1.Loc, &ListNode{Elts:$2}) } 
    | LBRACK eltlist SEMI RBRACK { $$ = newNode($1.Loc, &ListNode{Elts:$2}) } 

eltlist
    : exp { $$ = NewNodeList($1) }
    | eltlist SEMI exp { $$ = append($1, $3) }

tuple
    : LPAREN tuple_comps RPAREN
    { $$ = newNode($1.Loc, &TupleNode{Comps:$2}) } 

tuple_comps
    : exp COMMA exp { $$ = append(NewNodeList($1), $3) }
    | tuple_comps COMMA exp { $$ = append($1, $3) }

array
    : LBRACK PIPE PIPE RBRACK { $$ = newNode($1.Loc, &ArrayNode{}) } 
    | LBRACK PIPE eltarray PIPE RBRACK
      { $$ = newNode($1.Loc, &ArrayNode{Elts:$3}) } 
    | LBRACK PIPE eltarray SEMI PIPE RBRACK
      { $$ = newNode($1.Loc, &ArrayNode{Elts:$3}) } 

eltarray
    : exp { $$ = NewNodeList($1) }
    | eltarray SEMI exp { $$ = append($1, $3) }

binexp
    : exp ADD exp
    { $$ = newNode($1.unionLoc($3), &AddNode{Left:$1, Right:$3}) }
    | exp SUB exp
    { $$ = newNode($1.unionLoc($3), &SubNode{Left:$1, Right:$3}) }
    | exp MUL exp
    { $$ = newNode($1.unionLoc($3), &MulNode{Left:$1, Right:$3}) }
    | exp DIV exp
    { $$ = newNode($1.unionLoc($3), &DivNode{Left:$1, Right:$3}) }
    | exp MOD exp
    { $$ = newNode($1.unionLoc($3), &ModNode{Left:$1, Right:$3}) }
    | exp ADD_DOT exp
    { $$ = newNode($1.unionLoc($3), &FAddNode{Left:$1, Right:$3}) }
    | exp SUB_DOT exp
    { $$ = newNode($1.unionLoc($3), &FSubNode{Left:$1, Right:$3}) }
    | exp MUL_DOT exp
    { $$ = newNode($1.unionLoc($3), &FMulNode{Left:$1, Right:$3}) }
    | exp DIV_DOT exp
    { $$ = newNode($1.unionLoc($3), &FDivNode{Left:$1, Right:$3}) }
    | exp EQ exp
    { $$ = newNode($1.unionLoc($3), &EqNode{Left:$1, Right:$3}) }
    | exp NE exp
    { $$ = newNode($1.unionLoc($3), &NeNode{Left:$1, Right:$3}) }
    | exp LT exp
    { $$ = newNode($1.unionLoc($3), &LtNode{Left:$1, Right:$3}) }
    | exp LE exp
    { $$ = newNode($1.unionLoc($3), &LeNode{Left:$1, Right:$3}) }
    | exp GT exp
    { $$ = newNode($1.unionLoc($3), &GtNode{Left:$1, Right:$3}) }
    | exp GE exp
    { $$ = newNode($1.unionLoc($3), &GeNode{Left:$1, Right:$3}) }
    | exp BAND exp
    { $$ = newNode($1.unionLoc($3), &BandNode{Left:$1, Right:$3}) }
    | exp BOR exp
    { $$ = newNode($1.unionLoc($3), &BorNode{Left:$1, Right:$3}) }
    | exp BXOR exp
    { $$ = newNode($1.unionLoc($3), &BxorNode{Left:$1, Right:$3}) }
    | exp LSHIFT exp
    { $$ = newNode($1.unionLoc($3), &LshiftNode{Left:$1, Right:$3}) }
    | exp RSHIFT exp
    { $$ = newNode($1.unionLoc($3), &RshiftNode{Left:$1, Right:$3}) }
    | exp CONCAT exp
    { $$ = newNode($1.unionLoc($3), &ConcatNode{Left:$1, Right:$3}) }
    | exp DOL exp
    {
        args := make([]*Node, 1)
        args[0] = $3
        $$ = newNode($1.unionLoc($3), &AppNode{Exp:$1, Args:args})
    }

args
    : arg { $$ = NewNodeList($1) }
    | args arg { $$ = append($1, $2) }

arg
    : simple_exp { $$ = $1 }
    | LIDENT COLON simple_exp
    { $$ = newNode($1.Loc, &LabeledArgNode{Name:$1, Exp:$3}) }

typdef
    : TYPE typdefbody {}
    | TYPE typdefbody typdefbodylist {}

typdefbodylist
    : AND typdefbody {}
    | typdefbodylist AND typdefbody {}

typdefbody
    : typconstrname typinfo {}
    | typparamlist typconstrname typinfo {}

typinfo
    : typconstraintlist {}
    | typeq typconstraintlist {}
    | typeq typrepr typconstraintlist {}
    | typrepr typconstraintlist {}

typconstraintlist
    : typconstraint {}
    | typconstraintlist typconstraint {}

typeq
    : typexp { $$ = $1 }

typrepr
    : EQ constrdecl {}
    | EQ PIPE constrdecl {}
    | EQ PIPE constrdecl constrdecllist{}
    | EQ constrdecl constrdecllist {}
    | EQ LBRACE fielddecllist_semi RBRACE {}
    | EQ LBRACE fielddecllist_semi SEMI RBRACE {}

constrdecllist
    : PIPE constrdecl {}
    | constrdecllist PIPE constrdecl {}

fielddecllist_semi
    : fielddecl {}
    | fielddecllist_semi SEMI fielddecl {}

fielddecl
    : fieldname COLON polytypexp {}
    | MUTABLE fieldname COLON polytypexp {}

fieldname
    : LIDENT { $$ = $1 }

typparamlist
    : typparam {}
    | LPAREN typparam typparamlist_comma RPAREN

typparamlist_comma
    : COMMA typparam {}
    | typparamlist_comma COMMA typparam {}

typparam
    : QUOTE UIDENT {}
    | variance QUOTE UIDENT {}

variance
    : ADD {}
    | SUB {}

typconstraint
    : CONSTRAINT QUOTE UIDENT EQ typexp {}

polytypexp
    : typexp {}
    | polynamelist DOT typexp {}

polynamelist
    : QUOTE LIDENT
    | polynamelist QUOTE LIDENT {}

typexp
    : raw_typexp %prec prec_raw_typexp { $$ = WrapWithPolyNode($1) }
    | raw_typexp AS QUOTE LIDENT
    { $$ = newNode($1.Loc, &TypeAliasNode{Exp:$1, Name:$4}) }

raw_typexp
    : simple_typexp %prec prec_simple_typexp { $$ = $1 }
    | raw_typexp RARROW raw_typexp
    { $$ = newNode($1.unionLoc($3), &TypeArrowNode{Left:$1, Right:$3}) }
    | label raw_typexp RARROW raw_typexp
    { $$ = newNode($1.Loc, &TypeArrowNode{Label:$1, Left:$2, Right:$4}) }
    | raw_typexp typexplist_ast %prec prec_typexplist_ast
    { $$ = newNode($1.Loc, &TypeTupleNode{Comps:ConsNodeList($1, $2)}) }
    | raw_typexp typconstr
    { $$ = newNode($1.Loc, &TypeConstrAppNode{Exps:NewNodeList($1), Constr:$2}) }

simple_typexp
    : QUOTE LIDENT
    { $$ = newNode($1.Loc.Union($2.Loc), &TypeVarNode{Name:$2}) }
    | LPAREN typexp RPAREN { $$ = $2 }
    | typconstr { $$ = $1 }
    | LPAREN typexp RPAREN typconstr
    { $$ = newNode($1.Loc, &TypeConstrAppNode{Exps:NewNodeList($2), Constr:$4}) }
    | LPAREN typexp typexplist_comma RPAREN typconstr
    { $$ = newNode($1.Loc, &TypeParamConstrNode{Exps:ConsNodeList($2, $3), Constr:$5}) }

label
    : Q LIDENT COLON { $$ = newNode($1.Loc, &LabelNode{Opt:true, Name:$2}) }
    | LIDENT COLON { $$ = newNode($1.Loc, &LabelNode{Name:$1}) }

typexplist_ast
    : MUL typexp { $$ = NewNodeList($2) }
    | typexplist_ast MUL typexp { $$ = append($1, $3) }

typexplist_comma
    : COMMA simple_typexp { $$ = NewNodeList($2) }
    | typexplist_comma COMMA simple_typexp { $$ = append($1, $3) }

use
    : USE UIDENT
    { $$ = newNode($1.Loc, &UseNode{Trait:$2}) }
    | USE UIDENT LBRACE trait_params RBRACE
    { $$ = newNode($1.Loc, &UseNode{Trait:$2, Params:$4}) }
    | USE UIDENT LBRACE trait_params RBRACE WITH vallist
    { $$ = newNode($1.Loc, &UseNode{Trait:$2, Params:$4, Dir:TraitInclude, Vals:$7}) }
    | USE UIDENT LBRACE trait_params RBRACE WITHOUT vallist
    { $$ = newNode($1.Loc, &UseNode{Trait:$2, Params:$4, Dir:TraitExclude, Vals:$7}) }
    | USE UIDENT WITH vallist
    { $$ = newNode($1.Loc, &UseNode{Trait:$2, Dir:TraitInclude, Vals:$4}) }
    | USE UIDENT WITHOUT vallist
    { $$ = newNode($1.Loc, &UseNode{Trait:$2, Dir:TraitExclude, Vals:$4}) }

trait_params
    : LIDENT EQ typexp
    {
        $$ = make([]*Node, 1)
        $$[0] = newNode($1.Loc, &TraitParam{Name:$1, TypeExp:$3})
    }
    | trait_params LIDENT EQ typexp
    { $$ = append($1, newNode($2.Loc, &TraitParam{Name:$2, TypeExp:$4})) }

vallist
    : LIDENT
    {
        $$ = make([]*Word, 1)
        $$[0] = $1
    }
    | vallist LIDENT { $$ = append($1, $2) }

%%

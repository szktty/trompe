%{

open Ast
open Located

let create_binexp left op_loc op right =
  let op = create (Some op_loc) op in
  less @@ `Binexp (left, op, right)

let create_unexp op_loc op exp =
  let op = create (Some op_loc) op in
  less @@ `Unexp (op, exp)

%}

%token <Ast.text> UIDENT
%token <Ast.text> LIDENT
%token <Ast.text> DEREF_LIDENT
%token <Ast.text> CHAR
%token <Ast.text> STRING
%token <int Located.t> INT
%token <float Located.t> FLOAT
%token <Location.t> TRUE
%token <Location.t> FALSE
%token <Location.t> LPAREN
%token RPAREN
%token LBRACK
%token RBRACK
%token LBRACE
%token RBRACE
%token COMMA                        (* "," *)
%token DOT                          (* "." *)
%token DOT2                         (* ".." *)
%token COLON                        (* ":" *)
%token COLON2                       (* "::" *)
%token SEMI                         (* ";" *)
%token LARROW                       (* "<-" *)
%token RARROW                       (* "->" *)
%token BAR                          (* "|" *)
%token CARET                        (* "^" *)
%token <Location.t> AT              (* "@" *)
%token <Location.t> AMP             (* "&" *)
%token <Location.t> LPIPE           (* "<|" *)
%token <Location.t> RPIPE           (* "|>" *)
%token <Location.t> LCOMP           (* "<<" *)
%token <Location.t> RCOMP           (* ">>" *)
%token <Location.t> EQ              (* "=" *)
%token <Location.t> EQQ             (* "==" *)
%token <Location.t> NE              (* "!=" *)
%token <Location.t> LT              (* "<" *)
%token <Location.t> LE              (* "<=" *)
%token <Location.t> GT              (* ">" *)
%token <Location.t> GE              (* ">=" *)
%token <Location.t> PLUS            (* "+" *)
%token <Location.t> MINUS           (* "-" *)
%token <Location.t> AST             (* "*" *)
%token <Location.t> AST2            (* "**" *)
%token <Location.t> SLASH           (* "/" *)
%token <Location.t> PCT             (* "%" *)
%token <Location.t> AND             (* "and" *)
%token <Location.t> OR              (* "or" *)
%token DO                           (* "do" *)
%token CASE                         (* "case" *)
%token CATCH                        (* "catch" *)
%token DEF                          (* "def" *)
%token ELSE                         (* "else" *)
%token ELSEIF                       (* "elseif" *)
%token END                          (* "end" *)
%token FOR                          (* "for" *)
%token FUN                          (* "fun" *)
%token IF                           (* "if" *)
%token IN                           (* "in" *)
%token LET                          (* "let" *)
%token MODULE                       (* "module" *)
%token RAISE                        (* "raise" *)
%token RETURN                       (* "return" *)
%token PARTIAL                      (* "partial" *)
%token THEN                         (* "then" *)
%token TRY                          (* "try" *)
%token WHEN                         (* "when" *)
%token EOF

%nonassoc RAISE RETURN
%left OR AND
%right EQ LARROW
%nonassoc NE EQQ
%left LT GT LE GE
%left LCOMP RCOMP
%left RPIPE LPIPE
%left PLUS MINUS
%left AST SLASH PCT
%right AST2
%nonassoc prefix

%nonassoc app
%nonassoc LPAREN LBRACK

%start <Ast.t> prog

%%

prog:
  | block EOF { less @@ `Chunk $1 }

block:
  | (* empty *) { [] }
  | exp_list { $1 }

exp_list:
  | rev_exp_list { Core.Std.List.rev $1 }

rev_exp_list:
  | exp %prec prefix { [$1] }
  | rev_exp_list exp %prec prefix { $2 :: $1 }
  | rev_exp_list SEMI exp %prec prefix { $3 :: $1 }

exp:
  | module_def { $1 }
  | LET pattern EQ exp { Ast.nop }
  | LET LIDENT LARROW exp { less @@ `Refdef ($2, $4) }
  | var LARROW exp { less @@ `Assign ($1, $3) }
  | DO block END { less @@ `Block $2 }
  | RETURN exp { less @@ `Return $2 }
  | RAISE exp { less @@ `Raise $2 }
  | FOR LIDENT IN exp DO exp_list END
  { less @@ `For { for_var = $2; for_range = $4; for_block = $6 } }
  | fundef_exp { $1 }
  | if_exp { $1 }
  | case_exp { $1 }
  | fun_exp { $1 }
  | bin_exp { $1 }
  | unary_exp { $1 }
  | simple_exp { $1 }
  | directive { $1 }

module_def:
  | MODULE END { Ast.nop }
  | MODULE exp_list END { Ast.nop }

fundef_exp:
  | DEF LIDENT param_list EQ exp
  {
    less @@ `Fundef {
        fdef_name = $2;
        fdef_params = $3;
        fdef_block = [$5];
    }
  }
  | DEF LIDENT param_list block END
  {
    less @@ `Fundef {
        fdef_name = $2;
        fdef_params = $3;
        fdef_block = $4;
    }
  }

param_list:
  | LPAREN RPAREN { [] }
  | LPAREN param_list_body RPAREN { $2 }

param_list_body:
  | rev_param_list { Core.Std.List.rev $1 }

rev_param_list:
  | param { [$1] }
  | rev_param_list COMMA param { $3 :: $1 }

param:
  | LIDENT { $1 }
  | LIDENT COLON type_exp { $1 } (* TODO *)

if_exp:
  | IF exp THEN block END
  {
      less @@ `If {
          if_actions = [($2, $4)];
          if_else = [] }
  }
  | IF exp THEN block elseif_block END
  {
      less @@ `If {
          if_actions = ($2, $4) :: $5;
          if_else = [] }
  }
  | IF exp THEN block ELSE block END
  {
      less @@ `If {
          if_actions = [($2, $4)];
          if_else = $6 }
  }
  | IF exp THEN block elseif_block ELSE block END
  {
      less @@ `If {
          if_actions = ($2, $4) :: $5;
          if_else = $7 }
  }

elseif_block:
  | rev_elseif_block { Core.Std.List.rev $1 }

rev_elseif_block:
  | elseif_exp { [$1] }
  | rev_elseif_block elseif_exp { $2 :: $1 }

elseif_exp:
  | ELSEIF exp THEN block { ($2, $4) }

case_exp:
  | CASE exp DO case_clause_list END
  { less @@ `Case { case_val = $2; case_cls = $4 } }
  | CASE exp DO BAR case_clause_list END
  { less @@ `Case { case_val = $2; case_cls = $5 } }

case_clause_list:
  | rev_case_clause_list { Core.Std.List.rev $1 }

rev_case_clause_list:
  | case_clause { [$1] }
  | rev_case_clause_list BAR case_clause { $3 :: $1 }

case_clause:
  | case_pattern RARROW block
  { { Ast.case_cls_var = None;
      case_cls_ptn = fst $1;
      case_cls_guard = snd $1;
      case_cls_action = $3 } }
  | LIDENT EQ case_pattern RARROW block
  { { Ast.case_cls_var = Some $1;
      case_cls_ptn = fst $3;
      case_cls_guard = snd $3;
      case_cls_action = $5 } }

case_pattern:
  | pattern { ($1, None) }
  | pattern WHEN exp { ($1, Some $3) }

funcall_exp:
  | prefix_exp paren_arg_list
  {
    less @@ `Funcall { fc_fun = $1; fc_args = $2 }
  }

paren_arg_list:
  | LPAREN RPAREN { [] }
  | LPAREN arg_list RPAREN { $2 }

arg_list:
  | rev_arg_list { Core.Std.List.rev $1 }

rev_arg_list:
  | exp { [$1] }
  | rev_arg_list COMMA exp { $3 :: $1 }

fun_exp:
  | FUN fun_param_list RARROW block END { Ast.nop }

fun_param_list:
  | LPAREN RPAREN { [] }
  | LPAREN param_list_body RPAREN { $2 }
  | param_list_body { $1 }

bin_exp:
  | exp PLUS exp { create_binexp $1 $2 `Add $3 }
  | exp MINUS exp { create_binexp $1 $2 `Sub $3 }
  | exp AST exp { create_binexp $1 $2 `Mul $3 }
  | exp AST2 exp { create_binexp $1 $2 `Pow $3 }
  | exp SLASH exp { create_binexp $1 $2 `Div $3 }
  | exp PCT exp { create_binexp $1 $2 `Mod $3 }
  | exp EQQ exp { create_binexp $1 $2 `Eq $3 }
  | exp NE exp { create_binexp $1 $2 `Ne $3 }
  | exp LT exp { create_binexp $1 $2 `Lt $3 }
  | exp LE exp { create_binexp $1 $2 `Le $3 }
  | exp GT exp { create_binexp $1 $2 `Gt $3 }
  | exp GE exp { create_binexp $1 $2 `Ge $3 }
  | exp AND exp { create_binexp $1 $2 `And $3 }
  | exp OR exp { create_binexp $1 $2 `Or $3 }
  | exp LPIPE exp { create_binexp $1 $2 `Lpipe $3 }
  | exp RPIPE exp { create_binexp $1 $2 `Rpipe $3 }
  | exp LCOMP exp { create_binexp $1 $2 `Lcomp $3 }
  | exp RCOMP exp { create_binexp $1 $2 `Rcomp $3 }

unary_exp:
  | PLUS simple_exp { create_unexp $1 `Pos $2 }
  | MINUS simple_exp { create_unexp $1 `Neg $2 }
  | AST simple_exp { less @@ `Deref $2 }
  | DEREF_LIDENT { less @@ `Deref_var $1 } (* TODO: needed? *)

simple_exp:
  | prefix_exp %prec app { $1 }
  | literal { $1 }

prefix_exp:
  | var { $1 }
  | funcall_exp { $1 }
  | LPAREN exp RPAREN { $2 }
  | LPAREN exp COLON type_exp RPAREN { $2 } (* TODO *)

directive:
  | AT LIDENT paren_arg_list { less @@ `Directive ($2, $3) }

var:
  | LIDENT
  {
    create $1.loc @@ `Var {
        np_prefix = None;
        np_name = $1; }
  }
  | UIDENT
  {
    create $1.loc @@ `Path {
        np_prefix = None;
        np_name = $1; }
  }
  | prefix_exp DOT LIDENT
  { less @@ `Var {
        np_prefix = Some $1;
        np_name = $3; }
  }
  | prefix_exp DOT UIDENT
  { less @@ `Path {
        np_prefix = Some $1;
        np_name = $3; }
  }
  | prefix_exp LBRACK exp RBRACK
  { less @@ `Index { idx_prefix = $1; idx_index = $3 } }

literal:
  | LPAREN RPAREN { locate $1 `Unit }
  | STRING { create $1.loc @@ `String $1.desc }
  | INT { create $1.loc @@ `Int $1.desc }
  | FLOAT { create $1.loc @@ `Float $1.desc }
  | TRUE { locate $1 @@ `Bool true }
  | FALSE { locate $1 @@ `Bool false }
  | list_ { less @@ `List $1 }
  | tuple { less @@ `Tuple $1 }
  | INT DOT2 INT { less @@ `Range ($1, $3) }

list_:
  | LBRACK RBRACK { [] }
  | LBRACK elts RBRACK { $2 }

elts:
  | rev_elts { Core.Std.List.rev $1 }

rev_elts:
  | exp { [$1] }
  | rev_elts COMMA exp { $3 :: $1 }

tuple:
  | LPAREN exp COMMA rev_elts RPAREN { $2 :: (Core.Std.List.rev $4) }

pattern:
  | LPAREN pattern RPAREN { $2 }
  | LPAREN RPAREN { locate $1 @@ `Unit }
  | STRING { create $1.loc @@ `String $1.desc }
  | INT { create $1.loc @@ `Int $1.desc }
  | FLOAT { create $1.loc @@ `Float $1.desc }
  | TRUE { locate $1 @@ `Bool true }
  | FALSE { locate $1 @@ `Bool false }
  | LIDENT { create $1.loc @@ `Var $1 }
  | UIDENT { less @@ `Variant ($1, [])  }
  | UIDENT LPAREN elts_ptn RPAREN { less @@ `Variant ($1, $3)  }
  | CARET LIDENT { less @@ `Pin $2 }
  | pattern COLON2 LIDENT { less @@ `Cons ($1, $3) }
  | list_ptn { less @@ `List $1 }
  | tuple_ptn { less @@ `Tuple $1 }

list_ptn:
  | LBRACK RBRACK { [] }
  | LBRACK elts_ptn RBRACK { $2 }

elts_ptn:
  | rev_elts_ptn { Core.Std.List.rev $1 }

rev_elts_ptn:
  | pattern { [$1] }
  | rev_elts_ptn COMMA pattern { $3 :: $1 }

tuple_ptn:
  | LPAREN pattern COMMA rev_elts_ptn RPAREN { Core.Std.List.rev ($2 :: $4) }

type_exp:
  | simple_type_exp { $1 }
  | simple_type_exp simple_type_exp { $1 }
  | simple_type_exp LT type_exp_list GT { $1 }

type_exp_list:
  | rev_type_exp_list { Core.Std.List.rev $1 }

rev_type_exp_list:
  | type_exp { [] }
  | rev_type_exp_list COMMA type_exp { [] }

simple_type_exp:
  | LPAREN type_exp RPAREN { Ast.nop }
  | LIDENT { Ast.nop }
  | type_path { Ast.nop }
  | LBRACK type_exp RBRACK { Ast.nop }
  | LPAREN type_exp COMMA type_exp_list RPAREN { Ast.nop }

type_path:
  | rev_type_path { Core.Std.List.rev $1 }

rev_type_path:
  | UIDENT { [Ast.nop] } 
  | rev_type_path DOT UIDENT { Ast.nop :: $1 }

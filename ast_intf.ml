type op = op_desc Located.t

and op_desc = [
  | `Pos              (* "+" *)
  | `Neg              (* "-" *)
  | `Eq               (* "==" *)
  | `Ne               (* "!=" *)
  | `Lt               (* "<" *)
  | `Le               (* "<=" *)
  | `Gt               (* ">" *)
  | `Ge               (* ">=" *)
  | `List_add         (* "++" *)
  | `List_diff        (* "--" *)
  | `Add              (* "+" *)
  | `Sub              (* "-" *)
  | `Mul              (* "*" *)
  | `Div              (* "/" *)
  | `Pow              (* "**" *)
  | `Mod              (* "%" *)
  | `And              (* "and" *)
  | `Or               (* "or" *)
  | `Xor              (* "xor" *)
  | `Land             (* "band" *)
  | `Lor              (* "bor" *)
  | `Lxor             (* "bxor" *)
  | `Lcomp            (* "<<" *)
  | `Rcomp            (* ">>" *)
  | `Lpipe            (* "<|" *)
  | `Rpipe            (* "|>" *)
]

type t = desc Located.t

and desc = [
  | `Nop (* internal use *)
  | `Chunk of t list
  | `Refdef of text * t
  | `Assign of t * t
  | `Fundef of fundef
  | `Return of t
  | `Raise of t
  | `If of if_
  | `For of for_
  | `Case of case
  | `Block of t list
  | `Funcall of funcall
  | `Binexp of (t * op * t)
  | `Unexp of (op * t)
  | `Directive of (text * t list)
  | `Var of namepath
  | `Path of namepath
  | `Index of index
  | `Unit
  | `Bool of bool
  | `String of string
  | `Int of int
  | `Float of float
  | `List of t list
  | `Tuple of t list
  | `Range of (int Located.t * int Located.t)
  | `Enum of (text, t) enum
  | `Deref of t
  | `Deref_var of text
]

and fundef = {
  fdef_name : text;
  fdef_params : text list;
  fdef_block : t list;
}

and namepath = {
  np_prefix : t option;
  np_name : text;
}

and if_ = {
  if_actions : (t * t list) list;
  if_else : t list;
}

and for_ = {
  for_var : text;
  for_range : t;
  for_block : t list;
}

and case = {
  case_val : t;
  case_cls : case_cls list;
}

and case_cls = {
  case_cls_var : text option;
  case_cls_ptn : pattern;
  case_cls_guard : t option;
  case_cls_action : t list;
}

and funcall = {
  fc_fun : t;
  fc_args : t list;
}

and index = {
  idx_prefix : t;
  idx_index : t;
}

and ('name, 'value) enum = {
  enum_name : 'name;
  enum_params : ('name option * 'value);
}

and text = string Located.t

and pattern = ptn_desc Located.t

and ptn_desc = [
  | `Nop (* internal use *)
  | `Unit
  | `Bool of bool
  | `String of string
  | `Int of int
  | `Float of float
  | `Variant of (text * pattern list)
  | `Cons of (pattern * text)
  | `List of pattern list
  | `Tuple of pattern list
  | `Enum of (text, pattern) enum
  | `Var of text
  | `Pin of text
]

let nop : t = Located.less `Nop

let ptn_nop : pattern = Located.less `Nop

type op = op_desc Located.t

and op_desc = [
  | `Pos              (* "+" + integer *)
  | `Fpos             (* "+." + float *)
  | `Neg              (* "-" + integer *)
  | `Fneg             (* "-." + float *)
  | `Eq               (* "==" *)
  | `Ne               (* "!=" *)
  | `Lt               (* "<" *)
  | `Le               (* "<=" *)
  | `Gt               (* ">" *)
  | `Ge               (* ">=" *)
  | `List_add         (* "++" *)
  | `List_diff        (* "--" *)
  | `Add              (* "+" *)
  | `Fadd             (* "+." *)
  | `Sub              (* "-" *)
  | `Fsub             (* "-." *)
  | `Mul              (* "*" *)
  | `Fmul             (* "*." *)
  | `Div              (* "/" *)
  | `Fdiv             (* "/." *)
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
  | `Vardef of pattern * t
  | `Assign of assign
  | `Fundef of fundef
  | `Strdef of strdef
  | `Return of exp
  | `Raise of exp
  | `If of if_
  | `For of for_
  | `Case of case
  | `Block of exp_list
  | `Funcall of funcall
  | `Binexp of binexp
  | `Unexp of unexp
  | `Directive of (text * t list)
  | `Var of var
  | `Index of index
  | `Unit
  | `Bool of bool
  | `String of string
  | `Int of int
  | `Float of float
  | `List of exp_list
  | `Tuple of exp_list
  | `Range of (int Located.t * int Located.t)
  | `Struct of t struct_
  | `Enum of (text, t) enum
]

and exp = {
  exp : t;
  mutable exp_type : Type.t option;
}

and exp_list = {
  exp_list : t list;
  mutable exp_list_type : Type.t option;
}

and assign = {
  asg_var : t;
  asg_exp : t;
  mutable asg_type : Type.t option;
}

and fundef = {
  fdef_name : text;
  fdef_params : text list;
  fdef_block : t list;
  mutable fdef_param_types : Type.t list option;
  mutable fdef_type : Type.t option;
}

and strdef = {
  sdef_name : text;
  sdef_fields : sdef_field list;
  mutable sdef_type : Type.t option;
}

and sdef_field = {
  sdef_field_name : text;
  sdef_field_tyexp : t;
  mutable sdef_field_type : Type.t option;
}

and var = {
  var_prefix : t option;
  var_name : text;
  mutable var_type : Type.t option;
}

and if_ = {
  if_actions : (t * t list) list;
  if_else : t list;
  mutable if_type : Type.t option;
}

and for_ = {
  for_var : text;
  for_range : t;
  for_block : t list;
}

and case = {
  case_val : t;
  case_cls : case_cls list;
  mutable case_val_type : Type.t option;
  mutable case_cls_type : Type.t option;
}

and case_cls = {
  case_cls_var : text option;
  case_cls_ptn : pattern;
  case_cls_guard : t option;
  case_cls_action : t list;
  mutable case_cls_act_type : Type.t option;
}

and funcall = {
  fc_fun : t;
  fc_args : t list;
  mutable fc_fun_type : Type.t option;
  mutable fc_arg_types : Type.t list option;
}

and index = {
  idx_prefix : t;
  idx_index : t;
  mutable idx_type : Type.t option;
}

and 'a struct_ = {
  str_namepath : text Namepath.t;
  str_fields : (text * 'a option) list;
  mutable str_type : Type.t option;
}

and ('name, 'value) enum = {
  enum_name : 'name;
  enum_params : ('name option * 'value);
  mutable enum_type : Type.t option;
}

and binexp = {
  binexp_left : t;
  binexp_op : op;
  binexp_right : t;
  mutable binexp_type : Type.t option;
}

and unexp = {
  unexp_op : op;
  unexp_exp : t;
  mutable unexp_type : Type.t option;
}

and text = string Located.t

and pattern = {
  ptn_cls : ptn_cls Located.t;
  mutable ptn_type : Type.t option;
}

and ptn_cls = [
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

let ptn_nop : pattern =
  { ptn_cls = Located.less `Nop; ptn_type = None }

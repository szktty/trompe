(* 32bit word code *)

type length = int

type index = int

type name = string

type t =
  | Nop
  | Position of Position.t
  | Load_unit
  | Load_true
  | Load_false
  | Load_int of int
  | Load_some
  | Load_none
  | Load_literal of index (* index of value in literal list *)
  | Load_local of index
  | Store of index (* store pop local *)
  | Pop
  | Return
  | Label of name
  | Loop_head
  | Jump of index
  | Branch_true of index
  | Branch_false of index
  | Call of length
  | Primitive of index (* index of string as literal *)
  | List of length
  | Tuple of length

let id = function
  | Nop -> 0
  | Position _ -> 1
  | Load_unit -> 100
  | Load_true -> 101
  | Load_false -> 102
  | Load_int _ -> 103
  | Load_some -> 104
  | Load_none -> 105
  | Load_literal _ -> 106
  | Load_local _ -> 107
  | Store _ -> 200
  | Pop -> 300
  | Return -> 301
  | Label _ -> 302
  | Loop_head -> 303
  | Jump _ -> 304
  | Branch_true _ -> 305
  | Branch_false _ -> 306
  | Call _ -> 400
  | Primitive _ -> 401
  | List _ -> 500
  | Tuple _ -> 501

(* 32bit word code *)
type t =
  | Nop
  | Load_void
  | Load_true
  | Load_false
  | Load_int of int
  | Load_some
  | Load_none
  | Load_literal of int (* index of value in literal list *)
  | Load_local of int
  | Store of int (* store pop local *)
  | Pop
  | Return
  | Label of int
  | Loop_head
  | Jump of int
  | Branch_true of int
  | Branch_false of int
  | Primitive of int (* index of string as literal *)

let value = function
  | Nop -> 0
  | Load_void -> 100
  | Load_true -> 101
  | Load_false -> 102
  | Load_int -> 103
  | Load_some -> 104
  | Load_none -> 105
  | Load_literal -> 106
  | Load_local -> 107
  | Store -> 200
  | Pop -> 300
  | Return -> 301
  | Label -> 302
  | Loop_head -> 303
  | Jump -> 304
  | Branch_true -> 305
  | Branch_false -> 306
  | Primitive -> 400

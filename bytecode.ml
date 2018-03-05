(* Byte code file *)

open Base

type literal =
  | Int of int
  | String of string

type closure = {
  path : string list;
  name : string;
  lits : literal list;
  ops : Opcode.t list;
}

type t = {
  path : string list;
  name : string;
  closs : closure list;
}

let to_json (bc:t) =
  let to_str_list es =
    `List (List.map es ~f:(fun comp -> `String comp))
  in
  let to_lit = function
    | Int value -> `Int value
    | String value -> `String value
  in
  let to_op (op:Opcode.t) =
    let id = `Int (Opcode.id op) in
    match op with
    | Opcode.Nop
    | Load_unit
    | Load_true
    | Load_false
    | Load_some
    | Load_none
    | Pop
    | Return
    | Loop_head -> id

    | Load_int value
    | Load_literal value
    | Load_local value
    | Store value
    | Jump value
    | Branch_true value
    | Branch_false value
    | Call value
    | Primitive value 
    | List value
    | Tuple value -> `List [id; `Int value]

    | Label name -> `List [id; `String name]

    | Position pos -> `List [id; `Int pos.line; `Int pos.col]
  in
  let to_clos (clos:closure) =
    `Assoc [
      ("path", to_str_list clos.path);
      ("name", `String clos.name);
      ("literal_list", `List (List.map clos.lits ~f:to_lit));
      ("opcode_list", `List (List.map clos.ops ~f:to_op))
    ]
  in
  `Assoc [
    ("path", to_str_list bc.path);
    ("name", `String bc.name);
    ("closure_list", `List (List.map bc.closs ~f:to_clos))
  ]

let () =
  let path = ["foo"; "bar"] in
  let name = "baz" in
  let bc = {
    path;
    name;
    closs = [
      {
        path;
        name;
        lits = [String "show"; String "hello"];
        ops = [
          Nop;
          Primitive 0;
          Load_literal 0;
          Call 1;
        ]
      }
    ]
  }
  in
  Stdio.print_endline
    (Yojson.pretty_to_string (to_json bc) ~std:true)

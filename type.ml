open Core.Std

type t = desc Located.t

and desc = [
  | `Unit
  | `Bool
  | `Int
  | `Float
  | `String
  | `List of t
  | `Tuple of t list
  | `Range
  | `Fun of t list * t
               (*
  | `Module of module_
                *)
  | `Var of t option ref
]

and module_ = {
  mod_name : string;
  mod_attrs : t String.Map.t;
}

let create loc ty =
  Located.create loc ty

let create_tyvar loc =
  Located.create loc (`Var (ref None))

let rec to_string (ty:t) =
  match ty.desc with
  | `Unit -> "Unit"
  | `Bool -> "Bool"
  | `Int -> "Int"
  | `Float -> "Float"
  | `String -> "String"
  | `List e -> "[" ^ to_string e ^ "]"
  | `Tuple es ->
    "(" ^ String.concat ~sep:", " (List.map es ~f:to_string) ^ ")"
  | `Range -> "Range"
  | `Fun (params, ret) ->
    "((" ^ String.concat ~sep:", " (List.map params ~f:to_string) ^
    ") -> " ^ to_string ret ^ ")"
               (*
  | `Module of module_
                *)
  | `Var { contents = None } -> "?"
  | `Var { contents = Some ty } -> "?" ^ to_string ty


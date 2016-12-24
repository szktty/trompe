open Core.Std

type t = desc Located.t

and desc = [
  | `App of tycon * t list
  | `Var of tyvar
  | `Poly of tyvar list * t
  | `Meta of metavar
]

and tycon = [
  | `Unit
  | `Bool
  | `Int
  | `Float
  | `String
  | `List
  | `Tuple
  | `Range
  | `Option
  | `Struct of string list
  | `Enum of string list
  | `Fun
  | `Tyfun of tyvar list * t
  | `Unique of tycon * int
]

and tyvar = int

and metavar = t option ref

and module_ = t Module.t

let create loc ty =
  Located.create loc ty

let create_metavar loc =
  Located.create loc (`Meta (ref None))

let var_names = [|
  "a"; "b"; "c"; "d"; "e"; "f"; "g"; "h"; "i"; "j"; "k"; "l"; "m"; "n";
  "o"; "p"; "q"; "r"; "s"; "t"; "u"; "v"; "w"; "x"; "y"; "z";
|]

let rec to_string (ty:t) =
  match ty.desc with
  | `App (tycon, args) ->
    let tycon_s = match tycon with
      | `Unit -> "Unit"
      | `Bool -> "Bool"
      | `Int -> "Int"
      | `Float -> "Float"
      | `String -> "String"
      | `List -> "List"
      | `Tuple -> "Tuple"
      | `Range -> "Range"
      | `Fun -> "Fun"
      | `Option -> "Option"
      | _ -> failwith "not impl"
    in
    "App(" ^ tycon_s ^ ")"
  | `Meta { contents = None } -> "Meta(_)"
  | `Meta { contents = Some ty } -> "Meta(" ^ to_string ty ^ ")"
  | `Var n -> "Var(" ^ Array.get var_names n ^ ")"
  | `Poly (tyvars, ty) ->
    let names = List.map tyvars ~f:(Array.get var_names) in
    "Poly([" ^ (String.concat names ~sep:", ") ^ "], " ^ to_string ty ^ ")"

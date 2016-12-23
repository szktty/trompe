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
  | `Instance of int
]

and module_ = t Module.t

module Env = struct

  type env = {
    mutable refs : t option ref list;
  }

  let create () = { refs = [] }

  let instantiate env ref =
    match List.findi env.refs ~f:(fun _ ref' -> phys_equal ref ref') with
    | Some (n, _) -> `Instance n
    | None ->
      env.refs <- ref :: env.refs;
      `Instance (List.length env.refs)

end

let create loc ty =
  Located.create loc ty

let create_tyvar loc =
  Located.create loc (`Var (ref None))

let var_names = [|
  "a"; "b"; "c"; "d"; "e"; "f"; "g"; "h"; "i"; "j"; "k"; "l"; "m"; "n";
  "o"; "p"; "q"; "r"; "s"; "t"; "u"; "v"; "w"; "x"; "y"; "z";
|]

let rec to_string ?(debug=false) (ty:t) =
  let to_string = to_string ~debug in
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
  | `Var { contents = None } ->
    if debug then
      "?"
    else
      failwith "uninstantiated type"
  | `Var { contents = Some ty } ->
    if debug then
      "?" ^ to_string ty
    else
      to_string ty
  | `Instance n ->
    if debug then
      "?" ^ Int.to_string n
    else begin
      if Array.length var_names <= n then
        failwith ("too much type variables: " ^ Int.to_string n)
      else
        Array.get var_names n
    end

let to_repr ty = to_string ~debug:true ty

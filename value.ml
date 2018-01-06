open Base

type t =
  | Void
  | Int of int
  | String of string
  | List of t list
  | Tuple of t list
  | Prim of string
  | Clos of closure
  | Fun of t

and closure = {
  clos_env : t Map.M(String).t;
}

let string_exn = function
  | String s -> s
  | _ -> failwith "not string"

let list_exn = function
  | List es -> es
  | _ -> failwith "not list"


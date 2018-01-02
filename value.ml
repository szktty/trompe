open Base

type t =
  | Void
  | Int of int
  | String of string
  | Prim of string
  | Clos of closure
  | Fun of t
  | Path of string

and closure = {
  clos_env : t Map.M(String).t;
}

let string_exn = function
  | String s -> s
  | _ -> failwith "not string"

let path_exn = function
  | Path s -> s
  | _ -> failwith "not path"


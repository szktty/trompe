open Base

type t = {
  prims : prim Map.M(String).t
}

and context = {
  env : env;
}

and env = value Map.M(String).t

and value =
  | Void
  | Int
  | String
  | Prim of string
  | Clos of context * value

and prim = context -> value list -> (context, error) Result.t

and error = 
  | Invalid_arity
  | Invalid_type

let create () = {
  prims = Map.empty (module String);
}

let shared = ref (create ())

let add_prim ~name ~f =
  shared := { prims = Map.add !shared.prims ~key:name ~data:f }


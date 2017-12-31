open Base

type t = {
  rt_mods : module_ Map.M(String).t;
  rt_prims : prim Map.M(String).t;
}

and module_ = {
  mod_file : string option;
  mod_parent : module_;
  mod_ctx : context;
}

and context = {
  ctx_file : string option;
  ctx_env : env;
}

and env = value Map.M(String).t

and value =
  | Void
  | Int
  | String
  | Prim of string
  | Clos of context * value
  | Mod_fun of value

and prim = context -> value list -> (context, error) Result.t

and error = 
  | Invalid_arity
  | Invalid_type

let create () = {
  rt_mods = Map.empty (module String);
  rt_prims = Map.empty (module String);
}

let add_prim rt ~name ~f =
  { rt with rt_prims = Map.add rt.rt_prims ~key:name ~data:f }


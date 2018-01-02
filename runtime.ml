open Base

type t = {
  rt_mods : module_ Map.M(String).t;
  rt_prims : prim Map.M(String).t;
}

and module_ = {
  mod_file : string option;
  mod_name : string;
  mod_parent : module_;
  mod_subs : module_ Map.M(String).t;
  mod_ctx : context;
  mod_env : Value.t Map.M(String).t;
}

and context = {
  ctx_file : string option;
  ctx_env : Value.t Map.M(String).t;
  ctx_import : module_ list;
}

and args = {
  args_vals : Value.t list;
  args_tys : [`String] list;
}

and prim = {
  prim_fun : prim_fun;
  prim_arity : int;
}

and prim_fun = t -> context -> args -> (t * Value.t, error) Result.t

and error = 
  | Invalid_arity of int * int (* actual * expected *)
  | Invalid_type

let create m = {
  rt_mods = Map.empty (module String);
  rt_prims = Map.empty (module String);
}

let rec find_submod mods (path:Namepath.t) =
  match Map.find mods path.name with
  | None -> None
  | Some m ->
    match path.sub with
    | None -> Some m
    | Some sub -> find_submod m.mod_subs sub

let find_mod rt (path:Namepath.t) =
  find_submod rt.rt_mods path

    (*
let add_mod rt (path:Namepath.t) ~m =
     *)

let add_prim rt ~name ~f =
  { rt with rt_prims = Map.set rt.rt_prims ~key:name ~data:f }

module Args = struct

  let create values tys =
    { args_vals = values; args_tys = tys }

  let length (args:args) =
    List.length args.args_vals

  let value_exn args index =
    if length args <> (index + 1) then
      failwith "invalid index"
    else begin
      let arg = List.nth_exn args.args_vals index in
      let ty = List.nth_exn args.args_tys index in
      match arg, ty with
      | String _, `String -> arg
      | _ -> failwith "no arg"
    end

  let string_exn args index =
    match value_exn args index with
    | String s -> s
    | _ -> failwith "not string"

end


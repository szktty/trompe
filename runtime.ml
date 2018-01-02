open Base

type t = {
  rt_mods : module_ Map.M(String).t;
  rt_prims : prim Map.M(String).t;
}

and module_ = {
  mod_file : string option;
  mod_name : string;
  mod_parent : module_ option;
  mod_subs : module_ Map.M(String).t;
  mod_ctx : context option;
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
  prim_name : string;
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

let create_mod ?file ?parent name = 
  { mod_file = file;
    mod_name = name;
    mod_parent = parent;
    mod_subs = Map.empty (module String);
    mod_ctx = None;
    mod_env = Map.empty (module String);
  }

let find_mod rt path =
  With_return.with_return (fun r ->
      let mods =
        List.fold_left path ~init:rt.rt_mods
          ~f:(fun mods name ->
              match Map.find mods name with
              | None -> r.return None
              | Some m -> m.mod_subs)
      in
      Map.find mods (List.last_exn path))

let add_mod rt path ~m =
  let rec f mods path accu =
    match path with
    | [] -> Some (Map.set mods ~key:m.mod_name ~data:m)
    | name :: rest ->
      match Map.find mods name with
      | None -> None
      | Some m -> f m.mod_subs rest accu
  in
  match f rt.rt_mods path Map.empty with
  | None -> None
  | Some mods -> Some ({ rt with rt_mods = mods })

let add_prim rt ~name ~f ~arity =
  let prim = { prim_name = name;
               prim_fun = f;
               prim_arity = arity;
             } in
  { rt with rt_prims = Map.set rt.rt_prims ~key:name ~data:prim }

let add_prims rt (prims:(string * prim_fun * int) list) =
  List.fold_left prims ~init:rt
    ~f:(fun rt prim ->
        match prim with
        | name, f, arity -> add_prim rt ~name ~f ~arity)

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
    Value.string_exn (value_exn args index)

  let path_exn args index =
    Value.path_exn (value_exn args index)

end


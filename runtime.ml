open Base

type t = {
  rt_mods : module_ Map.M(String).t;
  rt_prims : prim Map.M(String).t;
}

and module_ = {
  mod_file : string option;
  mod_parent : string option;
  mod_name : string;
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

and prim_fun = t -> context -> args -> t * Value.t

let create m = {
  rt_mods = Map.empty (module String);
  rt_prims = Map.empty (module String);
}

let create_mod ?file path =
  { mod_file = file;
    mod_parent = Namepath.parent path;
    mod_name = Namepath.base path;
    mod_ctx = None;
    mod_env = Map.empty (module String);
  }

let find_mod rt path =
  Map.find rt.rt_mods path

let add_mod rt m =
  match m.mod_parent with
  | None ->
    Some({ rt with rt_mods = Map.set rt.rt_mods
                       ~key:m.mod_name
                       ~data:m})
  | Some parent ->
    match Map.find rt.rt_mods parent with
    | None -> None
    | Some _ ->
    let path = Namepath.concat parent m.mod_name in
    Some({ rt with rt_mods = Map.set rt.rt_mods
                       ~key:path
                       ~data:m})

let find_prim rt name =
  Map.find rt.rt_prims name

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

module Env = struct

  type t = Value.t Map.M(String).t

  let set env ~name ~value =
    Map.set env ~key:name ~data:value

end

module Module = struct

  type t = module_

  let set_attr m ~name ~value =
    { m with mod_env = Env.set m.mod_env ~name ~value }

end

let define rt ?(attrs=[]) ?(prims=[]) ~path () =
  let rt = add_prims rt prims in
  let m = create_mod path in
  let env = List.fold_left attrs ~init:m.mod_env
      ~f:(fun env attr ->
          match attr with
          | name, value ->
            begin match value with
              | Value.Prim prim ->
                begin match find_prim rt prim with
                  | Some _ -> ()
                  | None ->
                    failwith (Printf.sprintf "unknown primitive %s" prim)
                end
              | _ -> ()
            end;
            Env.set env ~name ~value) in
  let m = { m with mod_env = env } in
  Option.value_exn (add_mod rt m)

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

  let list_exn args index =
    Value.list_exn (value_exn args index)

end

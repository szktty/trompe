open Core.Std

type t = Lang.env

let debug env =
  let open Printf in
  printf "env{";
  String.Map.iteri Lang.(env.env_map)
    ~f:(fun ~key ~data ->
        printf "%s=%s; " key (Lang.to_string data));
  printf "}\n"

let import_module env m =
  Lang.(env.env_imports <- List.rev (m :: env.env_imports))

(* TODO: import したモジュールも探す *)
let rec find_value env name =
  (*printf "Env.find_value: %s from " name;debug env;*)
  match String.Map.find Lang.(env.env_map) name with
  | Some _ as res -> res
  | None ->
    match List.find_mapi env.env_imports
            ~f:(fun _ m -> find_value m.mod_env name) with
    | Some _ as v -> v
    | None ->
      match env.env_parent with
      | None -> None
      | Some env -> find_value env name

let add_value env name value =
  Lang.({ env with env_map = String.Map.add env.env_map ~key:name ~data:value })

let add_values env assocs =
  List.fold_left assocs ~init:env
    ~f:(fun env (key, value) -> add_value env key value)

(* target = env が属するモジュール *)
let create ?(imports=[]) ?target ?parent ?(values=[]) () =
  let env = { Lang.env_parent = parent;
              env_map = String.Map.empty;
              env_imports = imports;
              env_target = target }
  in
  add_values env values

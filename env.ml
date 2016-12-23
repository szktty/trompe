open Core.Std

type 'a t = {
  parent : 'a t option;
  attrs : 'a String.Map.t;
  imports : 'a t list;
}

let debug env ~f =
  let open Printf in
  printf "env{";
  String.Map.iteri env.attrs
    ~f:(fun ~key ~data ->
        printf "%s=%s; " key (f data));
  printf "}\n"

let import env m =
  { env with imports = List.rev (m :: env.imports) }

(* TODO: import したモジュールも探す *)
let rec find_attr env name =
  (*printf "Env.find_attr: %s from " name;debug env;*)
  match String.Map.find env.attrs name with
  | Some _ as res -> res
  | None ->
    match List.find_mapi env.imports
            ~f:(fun _ env -> find_attr env name) with
    | Some _ as v -> v
    | None ->
      match env.parent with
      | None -> None
      | Some env -> find_attr env name

let add_attr env name value =
  { env with attrs = String.Map.add env.attrs ~key:name ~data:value }

let add_attrs env attrs =
  List.fold_left attrs ~init:env
    ~f:(fun env (name, value) -> add_attr env name value)

let create ?(imports=[]) ?parent ?(attrs=[]) () =
  { parent = parent;
    attrs = String.Map.of_alist_reduce attrs ~f:(fun _ b -> b);
    imports = imports }

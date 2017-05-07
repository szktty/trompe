open Core.Std

type 'a t = {
  parent : 'a t option;
  name : string;
  mutable env : 'a Env.t;
  mutable submodules : 'a t list;
  mutable imports : 'a t list;
}

let create ?parent ?(submodules=[]) ?(imports=[]) ?env name =
  let env = Option.value env ~default:(Env.create ()) in
  { parent; name; env; submodules; imports }

let name m = m.name

let rec root m =
  match m.parent with
  | None -> Some m
  | Some m -> root m

let is_root m = Option.is_none m.parent

let import m x =
  m.imports <- x :: m.imports

let rec namepath m =
  match m.parent with
  | None -> Namepath.create m.name
  | Some m ->
    Namepath.create ~prefix:(Some (namepath m)) m.name

let rec find_module ?(prefix=[]) m ~name =
  match prefix with
  | [] -> failwith "must not be empty"
  | fst :: rest ->
    match find_module m ~name:fst with
    | None -> None
    | Some m -> find_module ~prefix:rest m ~name

let add_module m x =
  m.submodules <- x :: m.submodules

let rec find_attr m key =
  match Env.find m.env key with
  | Some _ as res -> res
  | None ->
    match List.find_mapi m.imports
            ~f:(fun _ m -> find_attr m key) with
    | Some _ as v -> v
    | None -> None

let add_attr m ~key ~data =
  m.env <- Env.add m.env ~key ~data

open Core.Std

type 'a t = {
  parent : 'a t option;
  imports : 'a Module.t list;
  attrs : 'a String.Map.t;
}

let create ?(parent=None) ?(imports=[]) ?attrs () =
  { parent; imports;
    attrs = Option.value attrs ~default:(String.Map.empty);
  }

let import env m =
  { env with imports = m :: env.imports }

let rec find env key =
  match String.Map.find env.attrs key with
  | Some _ as res -> res
  | None ->
    match List.find_mapi env.imports
            ~f:(fun _ m -> Module.find_attr m key) with
    | Some _ as res -> res
    | None ->
      match env.parent with
      | Some env -> find env key
      | None -> None

let add env ~key ~data =
  { env with attrs = String.Map.add env.attrs ~key ~data }

let concat env =
  let rec f attrs accu =
    String.Map.fold attrs ~init:accu
      ~f:(fun ~key ~data accu -> String.Map.add accu ~key ~data)
  in
  f env.attrs String.Map.empty

let merge env src =
  { env with attrs = String.Map.merge env.attrs src
                 ~f:(fun ~key owner ->
                     match owner with
                     | `Left v | `Right v -> Some v
                     | `Both (_, v2) -> Some v2)
  }

let debug env ~f =
  let print env indent =
    let open Printf in
    let indent_s = String.make (indent * 2) ' ' in
    printf "%s{\n" indent_s;
    String.Map.iteri env.attrs ~f:(fun ~key ~data ->
        printf "%s  %s = %s\n" indent_s key (f data));
    printf "%s}\n" indent_s
  in
  let rec f env indent =
    let indent = match env.parent with
      | None -> indent
      | Some parent ->
        f env indent;
        indent + 1
    in
    print env indent
  in
  f env 0

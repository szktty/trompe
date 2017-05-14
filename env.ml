open Core.Std

type 'a t = {
  imports : 'a Module.t list;
  attrs : 'a String.Map.t;
}

let create ?(imports=[]) ?attrs () =
  { imports;
    attrs = Option.value attrs ~default:(String.Map.empty);
  }

let import env m =
  { env with imports = m :: env.imports }

let rec find env key =
  match String.Map.find env.attrs key with
  | Some _ as res -> res
  | None -> List.find_mapi env.imports
              ~f:(fun _ m -> Module.find_attr m key)

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
  let open Printf in
  printf "{\n";
  String.Map.iteri env.attrs ~f:(fun ~key ~data ->
      printf "  %s = %s\n" key (f data));
  printf "}\n"

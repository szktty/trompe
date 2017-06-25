open Core.Std

type 'a t = {
  prefix : 'a t option;
  name : 'a;
}

let sep = "."

let create ?(prefix=None) name =
  { prefix; name }

let rec rev_names path =
  match path.prefix with
  | None -> [path.name]
  | Some path -> path.name :: rev_names path

let names path = List.rev @@ rev_names path

let rec iter path ~f =
  List.iter (names path) ~f

let to_string path =
  let rec f path =
    match path.prefix with
    | None -> [path.name]
    | Some path -> path.name :: f path
  in
  String.concat ~sep @@ names path

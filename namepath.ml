open Core.Std

type 'a t = {
  prefix : 'a t option;
  name : 'a;
}

let create ?(prefix=None) name =
  { prefix; name }

let rec to_rev_list path =
  match path.prefix with
  | None -> [path.name]
  | Some prefix -> path.name :: to_rev_list prefix

let to_list path =
  List.rev @@ to_rev_list path

let iter path ~f =
  List.iter (to_list path) ~f

let fold path ~init ~f =
  List.fold_right (to_rev_list path) ~init ~f

let to_string ?(sep=".") path =
  String.concat ~sep @@ to_list path

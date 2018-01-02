open Base

let defaults = [
  Lib_path.init;
  Lib_string.init;
]

let init rt =
  List.fold_left defaults ~init:rt ~f:(fun rt f -> f rt)

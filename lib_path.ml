open Base

let prim_new rt ctx args =
  Ok (rt, Value.Path (Runtime.Args.string_exn args 0))

let prim_string rt ctx args =
  Ok (rt, Value.String (Runtime.Args.path_exn args 0))

let init rt =
  Runtime.add_prims rt [
    ("path_new", prim_new, 1);
    ("path_string", prim_string, 1);
  ]

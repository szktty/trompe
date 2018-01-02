let prim_new rt ctx args =
  rt, Value.Path (Runtime.Args.string_exn args 0)

let prim_string rt ctx args =
  rt, Value.String (Runtime.Args.path_exn args 0)

let init rt =
  rt

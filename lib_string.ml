let prim_length rt ctx args =
  let s = Runtime.Args.string_exn args 0 in
  Ok (rt, Value.Int (String.length s))

let init rt =
  Runtime.add_prims rt [
    ("string_length", prim_length, 1)
  ]

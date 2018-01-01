let prim_length rt ctx args =
  let s = Runtime.Args.string_exn args 0 in
  rt, Runtime.Val_int (String.length s)

let init rt =
  rt

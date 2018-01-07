open Base
open Runtime

let prim_len rt ctx args =
  let s = Args.string_exn args 0 in
  (rt, Value.Int (String.length s))

let init rt =
  define rt
    ~name:"string"
    ~attrs:[
      ("len", Value.Prim "string_len");
    ]
    ~prims:[
      ("string_len", prim_len, 1)
    ]
    ()

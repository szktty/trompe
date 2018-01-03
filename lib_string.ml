open Base
open Runtime

let prim_length rt ctx args =
  let s = Args.string_exn args 0 in
  Ok (rt, Value.Int (String.length s))

let init rt =
  define rt
    ~name:"string"
    ~attrs:[
      ("length",Value.Prim "string_length");
    ]
    ~prims:[
      ("string_length", prim_length, 1)
    ]
    ()

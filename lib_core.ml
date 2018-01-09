open Base
open Runtime

module F = Caml.Filename

let prim_eq rt ctx args =
  let l = Args.value_exn args 0 in
  let r = Args.value_exn args 1 in
  rt, Value.Bool (Value.eq l r)

let init rt =
  define rt
    ~path:"core"
    ~attrs:[
      ("=", Value.Prim "core_eq");
    ]
    ~prims:[
      ("core_eq", prim_eq, 2);
    ]
    ()

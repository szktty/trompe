open Base
open Runtime

let prim_concat rt ctx args =
  let dir = Args.string_exn args 0 in
  let file = Args.string_exn args 1 in
  Ok (rt, Value.String (Caml.Filename.concat dir file))

let init rt =
  define rt
    ~name:"path"
    ~attrs:[
      ("concat", Value.Prim "path_concat");
    ]
    ~prims:[
      ("path_concat", prim_concat, 1);
    ]
    ()

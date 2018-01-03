open Base
open Runtime

let prim_new rt ctx args =
  Ok (rt, Value.Path (Args.string_exn args 0))

let prim_string rt ctx args =
  Ok (rt, Value.String (Args.path_exn args 0))

let init rt =
  define rt
    ~name:"path"
    ~attrs:[
      ("new", Value.Prim "path_new");
      ("string", Value.Prim "path_string");
    ]
    ~prims:[
      ("path_new", prim_new, 1);
      ("path_string", prim_string, 1);
    ]
    ()

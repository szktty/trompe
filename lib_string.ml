open Base
open Runtime

let prim_len rt ctx args =
  let s = Args.string_exn args 0 in
  (rt, Value.Int (String.length s))

let prim_eq rt ctx args =
  let s1 = Args.string_exn args 0 in
  let s2 = Args.string_exn args 1 in
  (rt, Value.Bool (String.equal s1 s2))

let prim_join rt ctx args =
  let es = Args.list_exn args 0 in
  let comps = List.map es ~f:Value.string_exn in
  let sep = Args.string_exn args 1 in
  rt, Value.String (String.concat comps ~sep)

let init rt =
  define rt
    ~name:"string"
    ~attrs:[
      ("len", Value.Prim "string_len");
      ("eq", Value.Prim "string_eq");
      ("join", Value.Prim "string_join");
    ]
    ~prims:[
      ("string_len", prim_len, 1);
      ("string_eq", prim_eq, 2);
      ("string_join", prim_eq, 2);
    ]
    ()

open Base
open Runtime

module F = Caml.Filename

let prim_dir rt ctx args =
  let path = Args.string_exn args 0 in
  Ok (rt, Value.String (F.dirname path))

let prim_base rt ctx args =
  let path = Args.string_exn args 0 in
  Ok (rt, Value.String (F.basename path))

let prim_ext rt ctx args =
  let path = Args.string_exn args 0 in
  Ok (rt, Value.String (F.extension path))

let prim_root rt ctx args =
  let path = Args.string_exn args 0 in
  Ok (rt, Value.String (F.remove_extension path))

let prim_join rt ctx args =
  let es = Args.list_exn args 0 in
  let comps = List.fold_left es
      ~init:[]
      ~f:(fun accu e -> (Value.string_exn e) :: accu)
  in
  match List.rev comps with
  | [] -> Ok (rt, Value.String "")
  | comp :: comps ->
    Ok (rt, Value.String (List.fold_left comps
                            ~init:comp
                            ~f:(fun path comp -> F.concat path comp)))

let prim_split rt ctx args =
  let path = Args.string_exn args 0 in
  let dir = F.dirname path in
  let base = F.basename path in
  Ok (rt, Value.Tuple [Value.String dir; Value.String base])

let prim_split_ext rt ctx args =
  let path = Args.string_exn args 0 in
  let root = F.remove_extension path in
  let ext = F.extension path in
  Ok (rt, Value.Tuple [Value.String root; Value.String ext])

let init rt =
  define rt
    ~name:"path"
    ~attrs:[
      ("sep", Value.String F.dir_sep);
      ("ext_sep", Value.String ".");
      ("dir", Value.Prim "path_dir");
      ("base", Value.Prim "path_base");
      ("ext", Value.Prim "path_ext");
      ("root", Value.Prim "path_root");
      ("join", Value.Prim "path_join");
      ("split", Value.Prim "path_split");
      ("split_ext", Value.Prim "path_split_ext");
    ]
    ~prims:[
      ("path_dir", prim_dir, 1);
      ("path_base", prim_base, 1);
      ("path_ext", prim_ext, 1);
      ("path_root", prim_root, 1);
      ("path_join", prim_join, 1);
      ("path_split", prim_split, 1);
      ("path_split_ext", prim_split_ext, 1);
    ]
    ()

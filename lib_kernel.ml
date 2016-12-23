open Core.Std
open Value

let prim_show ctx env args =
  match List.hd args with
  | None -> failwith "must be an argument"
  | Some arg ->
    Printf.printf "%s\n" (Value.to_string arg);
    `Unit

let install () =
  Primitive.register [
    ("show", prim_show);
  ];
  let env = Env.create ~values:[
      ("show", `Prim "show");
    ] ()
  in
  Module.(register @@ create "Kernel" env)

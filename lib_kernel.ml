open Core.Std
open Value

let prim_show args =
  match List.hd args with
  | None -> failwith "must be an argument"
  | Some arg ->
    Printf.printf "%s\n" (Value.to_string arg);
    `Unit

let install () =
  let primitives = [
    ("show", prim_show);
  ]
  in
  List.iter primitives
    ~f:(fun (name, primitive) -> Module.add_primitive ~name ~primitive);
  Module.define @@ Module.create ~name:"Kernel" ()

open Core.Std
open Value

let prim_show args =
  match List.hd args with
  | None -> failwith "must be an argument"
  | Some arg ->
    Printf.printf "%s\n" (Value.to_string arg);
    `Unit

let prim_printf args =
  Printf.printf "(not impl)\n";
  `Unit

let init () =
  Runtime.Spec.(define "Kernel"
                +> fun_ "show" Type.Spec.(a @-> unit) "show"
                +> fun_ "printf" Type.Spec.fun_printf "printf"
                +> string "version" "0.0.1"
                |> end_);
  Runtime.Primitive.add "show" prim_show;
  Runtime.Primitive.add "printf" prim_printf;
  ()

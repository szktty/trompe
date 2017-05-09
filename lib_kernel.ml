open Core.Std
open Value

let prim_show args =
  match List.hd args with
  | None -> failwith "must be an argument"
  | Some arg ->
    Printf.printf "%s\n" (Value.to_string arg);
    `Unit

let init () =
  Runtime.Spec.(define "Kernel"
                +> fun_ "show" Type.Spec.(a @-> unit) "show"
                +> string "version" "0.0.1"
                |> end_);
  Runtime.Primitive.add "show" prim_show

open Core.Std
open Value

let prim_to_string args : Value.t =
  Interp.Primitive.(
    match parse args [`Int] with
    | [`Int value] -> `String (Int.to_string value)
    | _ -> failwith "error")

let init () =
  Runtime.Spec.(define "list"
                +> fun_ "to_string" Type.Spec.(list a @-> string) "to_string"
                |> end_);
  Runtime.Primitive.add "to_string" prim_to_string

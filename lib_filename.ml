open Core.Std
open Value

let prim_to_string args : Value.t =
  Interp.Primitive.(
    match parse args [`Int] with
    | [`Int value] -> `String (Int.to_string value)
    | _ -> failwith "error")

let init () =
  Runtime.Spec.(define "filename"
                +> fun_ "to_string" Type.Spec.(int @-> string) "to_string"
                |> end_);
  Runtime.Primitive.add "to_string" prim_to_string

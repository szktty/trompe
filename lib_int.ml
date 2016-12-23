open Core.Std
open Value

let prim_to_string ctx env args : Value.t =
  Interp.Primitive.(
    match parse ctx args [`Int] with
    | [`Int value] -> `String (Int.to_string value)
    | _ -> failwith "error")

let install () =
  Primitive.register [
    ("int_to_string", prim_to_string)
  ];
  let env = Env.create ~attrs:[
      ("to_string", `Prim "int_to_string");
    ] ()
  in
  Interp.register @@ Module.create "Int" env

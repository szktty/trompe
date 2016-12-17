open Core.Std
open Lang

let prim_to_string ctx env args =
  Interp.Primitive.(
    match parse ctx args [`Int] with
    | [`Int value] -> env, `String (Int.to_string value)
    | _ -> failwith "error")

let install () =
  Primitive.register [
    ("int_to_string", prim_to_string)
  ];
  let env = Env.create ~values:[
      ("to_string", `Prim "int_to_string");
    ] ()
  in
  Module.(register @@ create "Int" env)

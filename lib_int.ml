open Core.Std
open Value

let prim_to_string args : Value.t =
  Interp.Primitive.(
    match parse args [`Int] with
    | [`Int value] -> `String (Int.to_string value)
    | _ -> failwith "error")

let install () = ()
                   (*
  let primitives = [
    ("int_to_string", prim_to_string)
  ]
  in
  List.iter primitives
    ~f:(fun (name, primitive) -> Module.add_primitive ~name ~primitive);
  Module.define @@ Module.create ~name:"Int" ()
                    *)

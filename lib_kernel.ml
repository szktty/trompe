open Core.Std
open Value

let prim_show args =
  match List.hd args with
  | None -> failwith "must be an argument"
  | Some arg ->
    Printf.printf "%s\n" (Value.to_string arg);
    `Unit

let prim_printf args =
  let buf = Buffer.create 16 in
  begin match args with
    | [] -> ()
    | `String fmt_s :: args ->
      let fmt = Utils.Format.parse fmt_s in
      let params = Utils.Format.params fmt in
      if List.length params <> List.length args then begin
        failwith "error"
      end;
      ignore @@ List.fold_left fmt
        ~init:args
        ~f:(fun args fmt ->
            match fmt with
            | Text c ->
              Buffer.add_char buf c;
              args
            | Int ->
              begin match args with
                | `Int v :: args ->
                  Buffer.add_string buf (Int.to_string v);
                  args
                | _ -> failwith "error"
              end
            | _ -> failwith "not impl")
    | _ -> ()
  end;
  Printf.printf "%s" (Buffer.contents buf);
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

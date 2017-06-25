open Core.Std
open Logging

let ext = ".go"

module Builder = struct

  type register = string

  type op =
    | Nop
    | Main of op list
    | Comment of string list
    | Package of string
    | Import of string
    | Struct of (string * type_) list
    | Fun of fun_
    | Var of (register * type_)
    | Return of register list
    | Defer of op
    | Get of register
    | Put of (register * register)
    | Call of (op * op list)
    | Terminal
    | Nil
    | Bool of bool
    | Int of int
    | Float of float
    | String of string

  and fun_ = {
    fname : string;
    fargs : (register * type_) list;
    fbody : op list;
    freturn : type_ list;
  }

  and type_ =
    | Ty_bool
    | Ty_int
    | Ty_float
    | Ty_string
    | Ty_error
    | Ty_ptr of type_
    | Ty_other of string

  type t = {
    mutable register : int;
    mutable ext : string String.Map.t;
    mutable funs : op list String.Map.t;
    mutable output : string option;
  }

  let assign builder =
    let id = builder.register + 1 in
    builder.register <- id;
    Printf.sprintf "r%d" id

  let create () =
    let bld = { register = 0;
                ext = String.Map.empty;
                funs = String.Map.empty;
                output = None }
    in
    bld.ext <- String.Map.add bld.ext ~key:"show" ~data:(assign bld);
    bld

  let add_fun bld name ops =
    bld.funs <- String.Map.add bld.funs ~key:name ~data:ops

  let rec write bld buf op =
    let open Buffer in
    let open Printf in
    match op with
    | Main ops ->
      Printf.printf "write program\n";
      add_string buf "package main\n\n";
      add_string buf "import \"fmt\"\n";
      add_string buf "func show(arg string) { fmt.Println(arg) }\n";

      String.Map.iteri bld.ext ~f:(fun ~key ~data ->
          add_string buf (sprintf "var %s = %s\n" data key));
      add_string buf "\n";

      add_string buf "func main() {\n";
      List.iter ops ~f:(write bld buf);
      add_string buf "}\n"

    | Call (f, args) -> 
      Printf.printf "call\n";
      let r1 = assign bld in
      add_string buf (r1 ^ " := ");
      write bld buf f;
      add_string buf "\n";
      let r2s = List.map args ~f:(fun arg ->
          let r2 = assign bld in
          add_string buf (r2 ^ " := ");
          write bld buf arg;
          add_string buf "\n";
          r2)
      in
      add_string buf (sprintf "%s(%s)\n" r1 (String.concat r2s ~sep:","));
      Printf.printf "finished -> %s\n" (Buffer.contents buf);

    | Var (reg, _) -> add_string buf reg

    | Int value -> Buffer.add_string buf @@ Int.to_string value
    | String value -> add_string buf @@ "\"" ^ value ^ "\""
    | _ -> Printf.printf "write nop\n";()

  let finish bld op_node =
    let buf = Buffer.create 1000 in
    write bld buf op_node;
    bld.output <- Some (Buffer.contents buf)

end

let rec compile' bld (env:string String.Map.t) node =
  let open Located in
  let open Builder in
  let open Ast in

  match node.desc with
  | `Chunk exps ->
    Printf.printf "chunk\n";
    let ops = snd @@ List.fold_left exps
        ~init:(env, [])
        ~f:(fun (env, accu) exp ->
            let (env', op) = compile' bld env exp in
            (env', op :: accu))
    in
    (env, Main (List.rev ops))

  | `Funcall fc ->
    Printf.printf "funcall\n";
    let f = easy_compile bld env fc.fc_fun in
    let args = List.map fc.fc_args
        ~f:(fun arg -> easy_compile bld env arg) in
    (env, Call (f, args))

  | `Var var ->
    let r1 = match String.Map.find env var.var_name.desc with
      | Some reg -> reg
      | None ->
        match String.Map.find bld.ext var.var_name.desc with
        | Some reg -> reg
        | None -> assign bld
    in
    (env, Var (r1, Ty_int)) (* TODO: type *)

  | `Int value -> (env, Int value)
  | `String value -> (env, String value)
  | _ -> (env, Nop)

and easy_compile bld env node = snd @@ compile' bld env node

let compile node ~file =
  let file =
    match String.rsplit2 file ~on:'.' with
    | None -> file ^ ext
    | Some (base, _) -> base ^ ext
  in
  verbosef "begin compile -> %s" file;
  let bld = Builder.create () in
  let op_node = snd @@ compile' bld String.Map.empty node in
  Builder.finish bld op_node;
  verbosef "end compile";
  Printf.printf "output -> %s\n" (Option.value_exn bld.Builder.output)

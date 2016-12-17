open Core.Std
open Located

include Ast_intf

let op_to_string op =
  match Located.(op.desc) with
  | `Eq -> "=="
  | `Ne -> "=="
  | `Lt -> "<"
  | `Le -> "<="
  | `Gt -> ">"
  | `Ge -> ">="
  | `Add -> "+"
  | `Sub -> "-"
  | `Mul -> "*"
  | `Div -> "/"
  | `Pow -> "**"
  | `Mod -> "%"
  | _ -> failwith "not supported operator"

let write_list chan es ~f =
  let open Out_channel in
  let open Printf in
  output_string chan "[";
  ignore @@ List.fold_left es
    ~f:(fun rest e ->
        f e;
        if rest > 0 then output_string chan " ";
        rest - 1)
    ~init:((List.length es) - 1);
  output_string chan "]"

let rec write chan node =
  let open Located in
  let open Out_channel in
  let open Printf in
  let output_string = output_string chan in
  let output_space () = output_string " " in
  let output_lp () = output_string "(" in
  let output_rp () = output_string ")" in
  let write = write chan in
  let write_op op = output_string @@ op_to_string op in
  let write_nodes es = write_list chan es ~f:write in
  let write_texts es = write_list chan es ~f:(fun e -> output_string e.desc) in
  (* TODO: location *)
  let output_namepath np =
    (* TODO: path *)
    output_string "\"";
    output_string np.np_name.desc;
    output_string "\""
  in
  match node.desc with
  | `Nop -> output_string "nop"
  | `Chunk exps ->
    output_string "(chunk ";
    write_nodes exps;
    output_string ")"
  | `Fundef fdef ->
    output_string "(fundef ";
    output_string fdef.fdef_name.desc;
    output_space ();
    write_texts fdef.fdef_params;
    output_space ();
    write_nodes fdef.fdef_block;
    output_string ")"
  | `Case case ->
    output_string "(case ";
    write case.case_val;
    output_string " [";
    write_list chan case.case_cls
      ~f:(fun cls ->
          output_lp ();
          write_ptn chan cls.case_cls_ptn;
          output_space ();
          begin match cls.case_cls_guard with
            | None -> output_string "true"
            | Some guard -> write guard
          end;
          output_space ();
          write_nodes cls.case_cls_action;
          output_rp ());
    output_string "]";
    output_rp ()
  | `Return exp ->
    output_string "(return ";
    write exp;
    output_rp ()
  | `If if_ ->
    output_string "(if [";
    write_list chan if_.if_actions
      ~f:(fun (cond, action) ->
          output_lp ();
          write cond;
          output_space ();
          write_nodes action;
          output_rp ());
    output_string "] ";
    write_nodes if_.if_else;
    output_rp ()
  | `For for_ ->
    output_string "(for ";
    output_string for_.for_var.desc;
    output_space ();
    write for_.for_range;
    output_space ();
    write_nodes for_.for_block;
    output_rp ()
  | `Funcall fc ->
    output_string "(funcall ";
    write fc.fc_fun;
    output_space ();
    write_nodes fc.fc_args;
    output_string ")"
  | `Binexp (left, op, right) ->
    output_string "(";
    write_op op;
    output_space ();
    write left;
    output_space ();
    write right;
    output_rp()
  | `Directive (name, args) ->
    output_string "(directive ";
    output_string name.desc;
    output_space ();
    write_nodes args;
    output_rp()
  | `Var np ->
    output_string "(var ";
    output_namepath np;
    output_rp ()
  | `Index idx ->
    output_string "(index ";
    write idx.idx_prefix;
    output_space ();
    write idx.idx_index;
    output_rp ()
  | `Unit -> output_string "()"
  | `Bool true -> output_string "true"
  | `Bool false -> output_string "false"
  | `Int v -> output_string @@ sprintf "%d" v
  | `Float v -> output_string @@ sprintf "%f" v
  | `String s -> output_string @@ sprintf "\"%s\"" s
  | `List exps ->
    output_string "(list ";
    write_nodes exps;
    output_rp ()
  | `Tuple exps ->
    output_string "(tuple ";
    write_nodes exps;
    output_rp ()
  | `Range (start, end_) ->
    output_string @@ sprintf "(range %d %d)" start.desc end_.desc
  | _ -> failwith "not supported"

and write_ptn chan ptn =
  let open Located in
  let open Out_channel in
  let open Printf in
  let output_string = output_string chan in
  let output_space () = output_string " " in
  let output_lp () = output_string "(" in
  let output_rp () = output_string ")" in
  let write = write_ptn chan in
  let write_ptns es = write_list chan es ~f:write in
  let write_texts es = write_list chan es ~f:(fun e -> output_string e.desc) in
  match ptn.desc with
  | `Unit -> output_string "()"
  | `Bool true -> output_string "true"
  | `Bool false -> output_string "false"
  | `Int v -> output_string @@ sprintf "%d" v
  | `Float v -> output_string @@ sprintf "%f" v
  | `String s -> output_string @@ sprintf "\"%s\"" s
  | `List ptns ->
    output_string "[";
    write_ptns ptns;
    output_string "]"
  | `Tuple ptns ->
    output_string "(tuple ";
    write_ptns ptns;
    output_rp ()
  | `Var name ->
    output_string @@ "(var " ^ name.desc ^ ")"
  | _ -> failwith "not supported pattern"

let print node =
  write Out_channel.stdout node;
  Printf.printf "\n"

let rec gen_loc (node:t) =
  let node = match node.desc with
    | `Nop
    | `Unit
    | `Bool _
    | `Int _
    | `Float _
    | `String _
    | `List _
    | `Tuple _ -> node
    | _ -> node
  in
  ignore @@ Option.try_with (fun () -> node.loc);
  node

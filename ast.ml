open Core.Std
open Located

include Ast_intf

let op_to_string op =
  match Located.(op.desc) with
  | `Pos -> "|+|"
  | `Fpos -> "|+.|"
  | `Neg -> "|-|"
  | `Fneg -> "|-.|"
  | `Eq -> "=="
  | `Ne -> "!="
  | `Lt -> "<"
  | `Le -> "<="
  | `Gt -> ">"
  | `Ge -> ">="
  | `Add -> "+"
  | `Fadd -> "+."
  | `Sub -> "-"
  | `Fsub -> "-."
  | `Mul -> "*"
  | `Fmul -> "*."
  | `Div -> "/"
  | `Fdiv -> "/."
  | `Pow -> "**"
  | `Mod -> "%"
  | _ -> failwith "not supported operator"

let location = function
  | _ -> Location.zero (* TODO *)

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

let rec write chan (node:Ast_intf.t) =
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
  let write_text text = output_string ("\"" ^ text.desc ^ "\"") in
  let write_texts es = write_list chan es ~f:write_text in
  match node with
  | `Nop _ -> output_string "nop"
  | `Chunk exps ->
    output_string "(chunk ";
    write_nodes exps;
    output_string ")"
  | `Vardef (ptn, exp) ->
    output_string "(vardef ";
    write_ptn chan ptn;
    output_space ();
    write exp;
    output_string ")"
  | `Fundef fdef ->
    output_string "(fundef ";
    output_string fdef.fdef_name.desc;
    output_space ();
    write_texts fdef.fdef_params;
    output_space ();
    write_nodes fdef.fdef_block;
    output_string ")"
  | `Strdef sdef ->
    output_string "(strdef ";
    write_text sdef.sdef_name;
    output_string " {";
    List.iter sdef.sdef_fields ~f:(fun fld ->
        write_text fld.sdef_field_name;
        output_string ":";
        write_tyexp chan fld.sdef_field_tyexp;
        output_string ", ");
    output_string "})"
  | `Assign assign ->
    output_string "(assign ";
    write assign.asg_var;
    output_space ();
    write assign.asg_exp;
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
    write exp.exp;
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
  | `Unexp exp ->
    output_string "(";
    write_op exp.unexp_op;
    output_space ();
    write exp.unexp_exp;
    output_string ")"
  | `Binexp exp ->
    output_string "(";
    write_op exp.binexp_op;
    output_space ();
    write exp.binexp_left;
    output_space ();
    write exp.binexp_right;
    output_rp ()
  | `Directive (name, args) ->
    output_string "(directive ";
    output_string name.desc;
    output_space ();
    write_nodes args;
    output_rp ()
  | `Var var ->
    output_string "(var ";
    Option.iter var.var_prefix ~f:(fun node ->
        write node;
        output_string " ");
    output_string "\"";
    output_string var.var_name.desc;
    output_string "\"";
    output_rp ()
  | `Index idx ->
    output_string "(index ";
    write idx.idx_prefix;
    output_space ();
    write idx.idx_index;
    output_rp ()
  | `Unit _ -> output_string "()"
  | `Bool { desc = true } -> output_string "true"
  | `Bool { desc = false } -> output_string "false"
  | `Int v -> output_string @@ sprintf "%d" v.desc
  | `Float v -> output_string @@ sprintf "%f" v.desc
  | `String s -> output_string @@ sprintf "\"%s\"" s.desc
  | `List exps ->
    output_string "(list ";
    write_nodes exps.exp_list;
    output_rp ()
  | `Tuple exps ->
    output_string "(tuple ";
    write_nodes exps.exp_list;
    output_rp ()
  | `Range (start, end_) ->
    output_string @@ sprintf "(range %d %d)" start.desc end_.desc
  | `Struct str ->
    output_string "{";
    Namepath.iter str.str_namepath ~f:(fun name ->
        write_text name;
        output_string ".");
    output_string ": ";
    List.iter str.str_fields ~f:(fun (name, v_opt) ->
        write_text name;
        Option.iter v_opt ~f:(fun v ->
            output_string "=";
            write v);
        output_string ", ");
    output_string "}"
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
  match ptn.ptn_cls with
  | `Unit _ -> output_string "()"
  | `Bool { desc = true } -> output_string "true"
  | `Bool { desc = false } -> output_string "false"
  | `Int v -> output_string @@ sprintf "%d" v.desc
  | `Float v -> output_string @@ sprintf "%f" v.desc
  | `String s -> output_string @@ sprintf "\"%s\"" s.desc
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

and write_tyexp chan tyexp =
  let open Located in
  let open Type in

  let write_list es =
    List.iter es ~f:(fun arg ->
        write_tyexp chan arg;
        output_string chan ", ")
  in

  match tyexp.desc with
  | Ty_var name ->
    output_string chan ("'" ^ name.desc)
  | Ty_namepath path ->
    Namepath.iter path ~f:(fun name -> output_string chan name.desc)
  | Ty_app (tycon, args) ->
    write_tyexp chan tycon;
    output_string chan "<";
    write_list args;
    output_string chan ">"
  | Ty_list e ->
    output_string chan "[";
    write_tyexp chan e;
    output_string chan "]"
  | Ty_tuple es ->
    output_string chan "(";
    write_list es;
    output_string chan ")"
  | Ty_option e ->
    write_tyexp chan e;
    output_string chan "?"

let print node =
  write Out_channel.stdout node;
  Printf.printf "\n"

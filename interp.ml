open Core.Std

module Context = Value.Context
module Exn = Value.Exn
module Op = Value.Op

module Error = struct

  type t = {
    context : Context.t;
    exn : Exn.t;
  }

  exception E of t

  let raise context exn =
    Pervasives.raise @@ E { context; exn }

end

let rec eval ctx env node =
  let open Ast in
  let open Context in
  let open Located in
  let open Printf in

  (* TODO: 環境を継続する exp_list とそうでない exp_list (リストリテラルなど) を別の関数にする *)
  let eval_exps ctx env exps =
    let (env, values) = List.fold_left exps ~init:(env, [])
        ~f:(fun (env, values) exp ->
            let (env, value) = eval ctx env exp in
            (env, value :: values))
    in
    (env, List.rev values)
  in

  (* TODO: これは環境を引き継ぐ *)
  let eval_block ctx env exps =
    match eval_exps ctx env exps with
    | (env, []) -> (env, `Unit)
    | (env, values) -> (env, List.last_exn values)
  in

  let eval_fundef ctx env value def_node def args =
    let env = Env.add env def.fdef_name.desc value in
    let env = List.fold2_exn def.fdef_params args ~init:env
        ~f:(fun env param arg -> Env.add env param.desc arg) in
    let ctx = Context.create ~parent:(Some ctx) ~callee:(Some def_node) () in
    let (_, values) = eval_exps ctx env def.fdef_block in
    values
  in

  match node.desc with
  | `Nop -> (env, `Unit)
  | `Chunk exps -> eval_block ctx env exps
  | `Unit -> (env, `Unit)
  | `Bool v -> (env, `Bool v)
  | `String s -> (env, `String s)
  | `Int v -> (env, `Int v)
  | `Float v -> (env, `Float v)
  | `Range (start, end_) -> (env, `Range (start.desc, end_.desc))

  | `List exps ->
    let (_, exps) = eval_exps ctx env exps.exp_list in
    (env, `List exps)

  | `Tuple exps ->
    let (_, exps) = eval_exps ctx env exps.exp_list in
    (env, `Tuple exps)

  | `Raise exp ->
    begin match eval ctx env exp.exp with
      | (_, `Exn e) -> Error.raise ctx (Exn.of_user_error e)
      | _ -> Error.raise ctx (Exn.of_reason Value_error "not exception")
    end

  | `Vardef (ptn, exp) ->
    let env, value = eval ctx env exp in
    begin match eval_ptn ctx env value ptn with
      | Result.Error () ->
        Error.raise ctx (Exn.of_reason Runtime_error "match failed")
      | Result.Ok env -> (env, value)
    end

  | `Fundef def ->
    let v = `Fun (def, Env.concat env) in
    let env = Env.add env ~key:def.fdef_name.desc ~data:v in
    (env, v)

  | `If if_ ->
    let cond = List.fold_left if_.if_actions ~init:None
        ~f:(fun ret (cond, action) ->
            match ret with
            | Some _ -> ret
            | None ->
              let (_, cond_val) = eval ctx env cond in
              match cond_val with
              | `Bool false -> None
              | `Bool true ->
                let (_, ret) = eval_block ctx env action in
                Some ret
              | _ -> failwith "if-exp condition must be bool")
    in
    begin match cond with
      | Some v -> (env, v)
      | None -> eval_block ctx env if_.if_else
    end

  | `For for_ ->
    let (op, start, end_) = match eval ctx env for_.for_range with
      | (_, `Range (start, end_)) ->
        if start <= end_ then
          ((+) 1, start, end_)
        else
          ((-) (-1), start, end_)
      | _ -> failwith "not range value"
    in
    let rec iter env start end_ =
      if start <= end_ then begin
        let env' = Env.add env for_.for_var.desc (`Int start) in
        ignore @@ eval_exps ctx env' for_.for_block;
        iter env' (op start) end_
      end
    in
    iter env start end_;
    (env, `Unit)

  | `Case case ->
    let (_, value) = eval ctx env case.case_val in
    let retval = List.find_mapi case.case_cls
        ~f:(fun _ cls ->
            let env = match cls.case_cls_var with
              | None -> env
              | Some var_ -> Env.add env var_.desc value
            in
            match eval_ptn ctx env value cls.case_cls_ptn with
            | Result.Error () -> None
            | Result.Ok env ->
              let guard = match cls.case_cls_guard with
                | None -> true
                | Some guard ->
                  match eval ctx env guard with
                  | (_, `Bool res) -> res
                  | _ -> failwith "guard must be bool"
              in
              if guard then begin
                Some (eval_block ctx env cls.case_cls_action)
              end else
                None)
    in
    begin match retval with
      | None -> failwith "pattern matching is not exhaustive"
      | Some (env, retval) -> (env, retval)
    end

  | `Funcall fc ->
    let validate_nargs nparams nargs =
      if nparams <> nargs then
        failwith (sprintf "given %d, expected %d" nargs nparams)
    in
    let (_, f) = eval ctx env fc.fc_fun in
    let (_, args) = eval_exps ctx env fc.fc_args in
    (* TODO: create new context *)
    begin match f with
      | `Prim name ->
        begin match Runtime.Primitive.find name with
          | None -> failwith ("unknown primitive: " ^ name)
          | Some f -> (env, f args)
        end
      | `Fun (def, capture) ->
        validate_nargs (List.length def.fdef_params) (List.length args);
        let fenv = Env.merge env capture in
        begin match eval_fundef ctx fenv f node def args with
          | [] -> (env, `Unit)
          | values -> (env, List.last_exn values)
        end
      | _ -> failwith (sprintf "%s is not function" (Value.to_string f))
    end

  | `Binexp exp ->
    let (_, left') = eval ctx env exp.binexp_left in
    let (_, right') = eval ctx env exp.binexp_right in
    let op = exp.binexp_op in
    let res = match op.desc with
      | `Eq -> Op.eq left' right'
      | `Ne -> Op.ne left' right'
      | `Le -> Op.le left' right'
      | `Add -> Op.add left' right'
      | `Sub -> Op.sub left' right'
      | `Mul -> Op.mul left' right'
      | `Div -> Op.div left' right'
      | `Mod -> Op.mod_ left' right'
      | _ -> failwith ("not yet supported operator: " ^ (Ast.op_to_string op))
    in
    (env, res)

  | `Directive (name, args) ->
    let (_, args) = eval_exps ctx env args in
    begin match name.desc with
      | "primitive" ->
        if List.length args = 0 then
          failwith "needs primitive name";
        let prim = match List.hd_exn args with
          | `String s -> s
          | v -> failwith ("primitive name must be string: " ^ (Value.to_string v))
        in
        let f = match Runtime.Primitive.find prim with
          | None -> failwith ("unknown primitive: " ^ prim)
          | Some f -> f
        in
        let fargs = match List.tl args with
          | None -> []
          | Some args -> args
        in
        (env, f fargs)
      | _ -> failwith ("unknown directive: " ^ name.desc)
    end;

  | `Var np ->
    (* TODO: get module from path *)
    begin match Env.find env np.np_name.desc with
      | None -> Error.raise ctx
                  (Exn.of_reason Name_error ("not found var: " ^ np.np_name.desc))
      (*failwith ("not found var: " ^ np.np_name.desc)*)
      | Some v -> (env, v)
    end

  | `Refdef (name, init) ->
    let (_, value) = eval ctx env init in
    let refval = `Ref (ref value) in
    let env = Env.add env name.desc refval in
    (env, refval)

  | `Assign asg ->
    begin match eval ctx env asg.asg_var with
      | (_, `Ref ref_) ->
        let (_, newval) = eval ctx env asg.asg_exp in
        ref_ := newval;
        (env, newval)
      | _ -> failwith "assign: not reference"
    end

  | `Deref exp ->
    begin match eval ctx env exp.exp with
      | (_, `Ref ref_) -> (env, !ref_)
      | _ -> failwith "deref: not reference"
    end

  | `Deref_var name ->
    begin match Env.find env name.desc with
      | None -> failwith ("not found var: " ^ name.desc)
      | Some v ->
        match v with
        | `Ref ref_ -> (env, !ref_)
        | _ -> failwith "deref var: not reference"
    end

  | _ ->
    Ast.write Out_channel.stdout node;
    Printf.printf "\n";
    failwith "not supported node"

and eval_ptn ctx env value ptn : (Value.t Env.t, unit) Result.t =
  let open Result in
  let test op env x y = if op x y then Ok env else Error () in
  match (ptn.ptn_cls.desc, value) with
  | (`Unit, `Unit) -> Ok env
  | (`Unit, _) -> Error ()
  | (`Bool true, `Bool true) -> Ok env
  | (`Bool false, `Bool false) -> Ok env
  | (`Bool _, _) -> Error ()
  | (`String x, `String y) -> test String.equal env x y
  | (`String _, _) -> Error ()
  | (`Int x, `Int y) -> test (=) env x y
  | (`Int _, _) -> Error ()
  | (`Float x, `Float y) -> test (=.) env x y
  | (`Float _, _) -> Error ()

  | (`Var name, _) ->
    if name.desc = "_" then
      Ok env
    else
      Ok (Env.add env ~key:name.desc ~data:value)

  | (`Pin name, _) ->
    begin match Env.find env name.desc with
      | None -> failwith "not found pin variable" (* TODO: exception *)
      | Some pin -> test Value.equal env pin value
    end

  | (`List xs, `List ys) when List.length xs = List.length ys ->
    List.fold2_exn xs ys
      ~init:(Ok env)
      ~f:(fun env x y ->
          match env with
          | Error () -> Error ()
          | Ok env -> eval_ptn ctx env y x)

  | (`List _, _) -> Error ()

  | (`Tuple xs, `Tuple ys) when List.length xs = List.length ys ->
    List.fold2_exn xs ys ~init:(Ok env)
      ~f:(fun env x y ->
          match env with
          | Error () -> Error ()
          | Ok env -> eval_ptn ctx env y x)

  | (`Tuple _, _) -> Error ()

  | _ -> failwith "eval pattern not impl"

let run node =
  let m = Module.create "__Main" in
  let ctx = Context.create ~belong:(Some m) () in
  let env = Runtime.value_env () in
  ignore @@ eval ctx env node

module Primitive = struct

  type arg = [
    | `Unit
    | `Bool
    | `String
    | `Int
    | `Float
    | `Range
    | `List of arg list
    | `Tuple of arg list
    | `Fun of arg list
    | `Ref of arg
    | `Exn
  ]

  (* TODO: deprecated *)
  let check_arg arg ty =
    match arg, ty with
    | `Int value, `Int -> arg
    | _ -> failwith "type error"

  (* TODO: deprecated *)
  let parse args types =
    if List.length args <> List.length types then begin
      List.map2_exn args types ~f:check_arg
    end else
      failwith "arity error"

  let check_arity name arity args =
    let count = List.length args in
    if count <> arity then
      failwith (sprintf "%s() takes %d arguments (%d given)" name arity count)

  let unit args i =
    match List.nth_exn args i with
    | `Unit -> ()
    | _ -> failwith (sprintf "argument %d must be unit" i)

  let get_bool args i =
    match List.nth_exn args i with
    | `Bool v -> v
    | _ -> failwith (sprintf "argument %d must be bool" i)

  let get_int args i =
    match List.nth_exn args i with
    | `Int v -> v
    | _ -> failwith (sprintf "argument %d must be int" i)

  let get_float args i =
    match List.nth_exn args i with
    | `Float v -> v
    | _ -> failwith (sprintf "argument %d must be float" i)

  let get_string args i =
    match List.nth_exn args i with
    | `String s -> s
    | _ -> failwith (sprintf "argument %d must be string" i)

  let get_range args i =
    match List.nth_exn args i with
    | `Range v -> v
    | _ -> failwith (sprintf "argument %d must be range" i)

  let get_list args i =
    match List.nth_exn args i with
    | `List v -> v
    | _ -> failwith (sprintf "argument %d must be list" i)

  let get_tuple args i =
    match List.nth_exn args i with
    | `Tuple v -> v
    | _ -> failwith (sprintf "argument %d must be tuple" i)

  let get_stream args i =
    match List.nth_exn args i with
    | `Stream v -> v
    | _ -> failwith (sprintf "argument %d must be stream" i)

end

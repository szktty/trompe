open Core.Std
open Located
open Type

type unity_exn = {
  uniexn_ex : Type.t;
  uniexn_ac : Type.t;
}

type mismatch = {
  mismatch_e: Ast.t;
  mismatch_ex : Type.t;
  mismatch_ac : Type.t;
}

exception Unify_error of unity_exn
exception Type_mismatch of mismatch
exception Deref_error of Type.t * string

(* for pretty printing (and type normalization) *)
(* 型変数を中身でおきかえる関数 (caml2html: typing_deref) *)
let rec deref_type (ty:Type.t) : Type.t =
  let create = Located.create ty.loc in
  match ty.desc with
      (*
  | `Fun(t1s, t2) -> `Fun(List.map deref_typ t1s, deref_typ t2)
  | `Tuple(ts) -> `Tuple(List.map deref_typ ts)
  | `Array(t) -> `Array(deref_typ t)
       *)
  | `List ty -> create @@ `List (deref_type ty)
  | `Var { contents = None } ->
    (* TODO: ここを a, b, c, ... にしたい
     * Map か何かに ref を登録する
    *)
    raise @@ Deref_error (ty, "uninstantiated type variable")
  | `Var ({ contents = Some ty } as ref) ->
    let ty' = deref_type ty in
    ref := Some ty';
    ty'
  | _ -> ty

let rec deref_id_type x ty = (x, deref_type ty)

let rec deref_term (e:Ast.t) : Ast.t =
  let map = List.map ~f:deref_term in
  let desc = match e.desc with
    | `Nop
    | `Unit
    | `Bool _
    | `String _
    | `Int _
    | `Float _
    | `Range _ -> e.desc

    | `Chunk es -> `Chunk (map es)
    | `Return e -> `Return (deref_term e)
    | `Fundef fdef ->
      `Fundef { fdef with
                fdef_block = List.map fdef.fdef_block  ~f:deref_term }
    | `Var path ->
      begin match path.np_prefix with
        | None -> `Var { path with np_type = deref_type path.np_type }
        | Some _ -> failwith "not yet supported"
      end

    | `List es -> `List (map es)
    | `Tuple es -> `Tuple (map es)
      (*
  | Not(e) -> Not(deref_term e)
  | Neg(e) -> Neg(deref_term e)
  | Add(e1, e2) -> Add(deref_term e1, deref_term e2)
  | Sub(e1, e2) -> Sub(deref_term e1, deref_term e2)
  | Eq(e1, e2) -> Eq(deref_term e1, deref_term e2)
  | LE(e1, e2) -> LE(deref_term e1, deref_term e2)
  | FNeg(e) -> FNeg(deref_term e)
  | FAdd(e1, e2) -> FAdd(deref_term e1, deref_term e2)
  | FSub(e1, e2) -> FSub(deref_term e1, deref_term e2)
  | FMul(e1, e2) -> FMul(deref_term e1, deref_term e2)
  | FDiv(e1, e2) -> FDiv(deref_term e1, deref_term e2)
  | If(e1, e2, e3) -> If(deref_term e1, deref_term e2, deref_term e3)
  | Let(xt, e1, e2) -> Let(deref_id_typ xt, deref_term e1, deref_term e2)
  | LetRec({ name = xt; args = yts; body = e1 }, e2) ->
    LetRec({ name = deref_id_typ xt;
             args = List.map deref_id_typ yts;
             body = deref_term e1 },
           deref_term e2)
  | App(e, es) -> App(deref_term e, List.map deref_term es)
  | Tuple(es) -> Tuple(List.map deref_term es)
  | LetTuple(xts, e1, e2) -> LetTuple(List.map deref_id_typ xts, deref_term e1, deref_term e2)
  | Array(e1, e2) -> Array(deref_term e1, deref_term e2)
  | Get(e1, e2) -> Get(deref_term e1, deref_term e2)
  | Put(e1, e2, e3) -> Put(deref_term e1, deref_term e2, deref_term e3)
       *)
    | _ -> failwith "TODO: impl"
  in
  create e.loc desc

let rec occur (ref:t option ref) (ty:Type.t) : bool =
  match ty.desc with
  | `Fun (params, ret) ->
    List.exists params ~f:(occur ref) || occur ref ret
  | `List ty -> occur ref ty
  (*| `Tuple tys -> List.exists tys ~f:(occur ref) *)
  | `Var ref2 when phys_equal ref ref2 -> true
  | `Var { contents = None } -> false
  | `Var { contents = Some t2 } -> occur ref t2
  | _ -> false

(* 型が合うように、型変数への代入をする (caml2html: typing_unify) *)
let rec unify ~(ex:Type.t) ~(ac:Type.t) : unit =
  match ex.desc, ac.desc with
  | `Unit, `Unit
  | `Bool, `Bool
  | `Int, `Int
  | `Float, `Float
  | `String, `String -> ()
  | `List ex, `List ac -> unify ~ex ~ac
  | `Tuple ty1s, `Tuple ty2s when List.length ty1s = List.length ty2s ->
    List.iter2_exn ty1s ty2s ~f:(fun ty1 ty2 -> unify ~ex:ty1 ~ac:ty2)
  | `Fun (param1s, ret1), `Fun (param2s, ret2)
    when List.length param1s = List.length param2s ->
    List.iter2_exn param1s param2s
      ~f:(fun param1 param2 -> unify ~ex:param1 ~ac:param2);
    unify ~ex:ret1 ~ac:ret2

      (*
  | `Fun(t1s, t1'), `Fun(t2s, t2') ->
    (try List.iter2 unify t1s t2s
     with Invalid_argument("List.iter2") -> raise (Unify_error(t1, t2)));
    unify t1' t2'
  | `Tuple(t1s), `Tuple(t2s) ->
    (try List.iter2 unify t1s t2s
     with Invalid_argument("List.iter2") -> raise (Unify_error(t1, t2)))
  | `Array(t1), `Array(t2) -> unify t1 t2
       *)
  | `Var ex, `Var ac when phys_equal ex ac -> ()
                                        (*
  | `Var({ contents = Some(t1') }), _ -> unify t1' t2
  | _, `Var({ contents = Some(t2') }) -> unify t1 t2'
  | `Var({ contents = None } as r1), _ -> (* 一方が未定義の型変数の場合 (caml2html: typing_undef) *)
    if occur r1 t2 then raise (Unify_error(t1, t2));
    r1 := Some(t2)
  | _, `Var({ contents = None } as r2) ->
    if occur r2 t1 then raise (Unify_error(t1, t2));
    r2 := Some(t1)
                                         *)
  | _, _ ->
    raise (Unify_error { uniexn_ex = ex; uniexn_ac = ac })

(* 型推論ルーチン (caml2html: typing_g) *)
let rec infer env (e:Ast.t) : (Type.t String.Map.t * Type.t) =
  Printf.printf "infer e: ";
  Ast.print e;
  let unit = Located.create e.loc `Unit in
  let easy_infer env e = snd @@ infer env e in
  let infer_block env es =
    List.fold_left es ~init:(env, Type.create e.loc `Unit)
      ~f:(fun (env, _) e -> infer env e)
  in
  try
    let env, ty = match e.desc with
      | `Nop -> (env, `Unit)

      | `Chunk es -> 
        let env, ty = List.fold_left es
            ~init:(env, unit)
            ~f:(fun (env, ty) e -> infer env e)
        in
        (env, ty.desc)

      | `Return e -> (env, (easy_infer env e).desc)

      | `Unit -> (env, `Unit)
      | `Bool _ -> (env, `Bool)
      | `Int _ -> (env, `Int)
      | `Float _ -> (env, `Float)
      | `String _ -> (env, `String)
      | `Range _ -> (env, `Range)

      | `List es ->
        begin match es with
          | [] -> (env, `List (Type.create_tyvar e.loc))
          | e :: es ->
            let base_ty = easy_infer env e in
            List.iter es ~f:(fun e ->
                unify ~ex:base_ty ~ac:(easy_infer env e));
            (env, `List base_ty)
        end

      | `Tuple es ->
        (env, `Tuple (List.map es ~f:(easy_infer env)))

      | `Fundef fdef ->
        let params = List.map fdef.fdef_params
            ~f:(fun param -> Type.create_tyvar param.loc)
        in
        let ret = Type.create_tyvar e.loc in
        let ty_desc = `Fun (params, ret) in
        let ty = Type.create e.loc ty_desc in
        let env = String.Map.add env ~key:fdef.fdef_name.desc ~data:ty in
        let fenv = List.fold2_exn fdef.fdef_params params ~init:env
            ~f:(fun env name ty ->
                String.Map.add env ~key:name.desc ~data:ty)
        in
        let _, ret' = infer_block fenv fdef.fdef_block in
        begin match ret.desc with
          | `Var { contents = Some ret } -> unify ~ex:ret ~ac:ret'
          | `Var ({ contents = None } as ref) -> ref := Some ret'
          | _ -> failwith "return type must be type variable"
        end;
        (env, ty_desc)

      | `Var path ->
        begin match Ast.(path.np_prefix) with
          | Some _ -> failwith "not yet supported"
          | None -> (env, path.np_type.desc)
        end
                (*
    | Not(e) ->
      unify `Bool (infer env e);
      `Bool
    | Neg(e) ->
      unify `Int (infer env e);
      `Int
    | Add(e1, e2) | Sub(e1, e2) -> (* 足し算（と引き算）の型推論 (caml2html: typing_add) *)
      unify `Int (infer env e1);
      unify `Int (infer env e2);
      `Int
    | FNeg(e) ->
      unify `Float (infer env e);
      `Float
    | FAdd(e1, e2) | FSub(e1, e2) | FMul(e1, e2) | FDiv(e1, e2) ->
      unify `Float (infer env e1);
      unify `Float (infer env e2);
      `Float
    | Eq(e1, e2) | LE(e1, e2) ->
      unify (infer env e1) (infer env e2);
      `Bool
    | If(e1, e2, e3) ->
      unify (infer env e1) `Bool;
      let t2 = infer env e2 in
      let t3 = infer env e3 in
      unify t2 t3;
      t2
    | Let((x, t), e1, e2) -> (* letの型推論 (caml2html: typing_let) *)
      unify t (infer env e1);
      g (M.add x t env) e2
    | Var(x) when M.mem x env -> M.find x env (* 変数の型推論 (caml2html: typing_var) *)
    | Var(x) when M.mem x !extenv -> M.find x !extenv
    | Var(x) -> (* 外部変数の型推論 (caml2html: typing_extvar) *)
      Format.eprintf "free variable %s assumed as external@." x;
      let t = `gentyp () in
      extenv := M.add x t !extenv;
      t
    | LetRec({ name = (x, t); args = yts; body = e1 }, e2) -> (* let recの型推論 (caml2html: typing_letrec) *)
      let env = M.add x t env in
      unify t (`Fun(List.map snd yts, g (M.add_list yts env) e1));
      infer env e2
    | App(e, es) -> (* 関数適用の型推論 (caml2html: typing_app) *)
      let t = `gentyp () in
      unify (infer env e) (`Fun(List.map (infer env) es, t));
      t
    | Tuple(es) -> `Tuple(List.map (infer env) es)
    | LetTuple(xts, e1, e2) ->
      unify (`Tuple(List.map snd xts)) (infer env e1);
      g (M.add_list xts env) e2
    | Array(e1, e2) -> (* must be a primitive for "polymorphic" typing *)
      unify (infer env e1) `Int;
      `Array(infer env e2)
    | Get(e1, e2) ->
      let t = `gentyp () in
      unify (`Array(t)) (infer env e1);
      unify `Int (infer env e2);
      t
    | Put(e1, e2, e3) ->
      let t = infer env e3 in
      unify (`Array(t)) (infer env e1);
      unify `Int (infer env e2);
      `Unit
       *)
      | _ -> Ast.print e; failwith "TODO"
    in
    let ty = Located.create e.loc ty in
    Printf.printf "inferred type: %s\n" (Type.to_string ty);
    (env, ty)
  with
  | Unify_error { uniexn_ex = ex; uniexn_ac = ac } ->
    raise @@ Type_mismatch {
      mismatch_e = deref_term e;
      mismatch_ex = deref_type ex;
      mismatch_ac = deref_type ac;
    }

let run (e:Ast.t) : Ast.t =
  ignore @@ infer String.Map.empty e;
  deref_term e

open Core.Std
open Located
open Type

type unity_exn = {
  uniexn_ex : Type.t;
  uniexn_ac : Type.t;
}

type mismatch = {
  mismatch_node: Ast.t;
  mismatch_ex : Type.t;
  mismatch_ac : Type.t;
}

exception Unify_error of unity_exn
exception Type_mismatch of mismatch
exception Deref_error of Type.t * string

let top_modules : Type.module_ String.Map.t ref = ref String.Map.empty

let register (m:Type.module_) =
  top_modules := String.Map.add !top_modules ~key:m.name ~data:m

let find_module name =
  String.Map.find !top_modules name

(* for pretty printing (and type normalization) *)
(* ���ѿ�����ȤǤ���������ؿ� (caml2html: typing_deref) *)
let rec deref_type env (ty:Type.t) : Type.t =
  let create = Located.create ty.loc in
  match ty.desc with
  | `App (tycon, args) ->
    create @@ `App (tycon, List.map args ~f:(deref_type env))
  | `Meta { contents = None } ->
    failwith "uninstantiated"
  | `Meta ({ contents = Some ty } as ref) ->
    let ty' = deref_type env ty in
    ref := Some ty';
    ty'
  | `Poly (metavars, ty) ->
    create @@ `Poly (metavars, deref_type env ty)
  | `Var _ -> ty

and deref_tycon env = function
      (*
  | `Fun(t1s, t2) -> `Fun(List.map deref_typ t1s, deref_typ t2)
  | `Tuple(ts) -> `Tuple(List.map deref_typ ts)
  | `Array(t) -> `Array(deref_typ t)
       *)
  | tycon -> tycon

let rec deref_id_type x ty = (x, deref_type ty)

let rec deref_term env (e:Ast.t) : Ast.t =
  let map = List.map ~f:(deref_term env) in
  let desc = match e.desc with
    | `Nop
    | `Unit
    | `Bool _
    | `String _
    | `Int _
    | `Float _
    | `Range _ -> e.desc

    | `Chunk es -> `Chunk (map es)
    | `Return e -> `Return (deref_term env e)

    | `Fundef fdef ->
      `Fundef { fdef with fdef_block = map fdef.fdef_block }

    | `Funcall call ->
      `Funcall { fc_fun = deref_term env call.fc_fun;
                 fc_args = map call.fc_args }

    | `Var path ->
      begin match path.np_prefix with
        | None -> `Var { path with np_type = deref_type env path.np_type }
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
  | `App (tycon, args) ->
    begin match occur_tycon ref tycon with
      | false -> false
      | true -> List.exists args ~f:(occur ref)
    end
  | `Meta ref2 when phys_equal ref ref2 -> true
  | `Meta { contents = None } -> false
  | `Meta { contents = Some t2 } -> occur ref t2
  | _ -> failwith "not impl"

and occur_tycon ref = function
  | `List | `Fun -> true
  | `Tyfun (_, ty) -> occur ref ty
  | _ -> failwith "not impl"

(* �����礦�褦�ˡ����ѿ��ؤ������򤹤� (caml2html: typing_unify) *)
let rec unify ~(ex:Type.t) ~(ac:Type.t) : unit =
  match ex.desc, ac.desc with
  | `App (`Unit, []), `App (`Unit, [])
  | `App (`Bool, []), `App (`Bool, [])
  | `App (`Int, []), `App (`Int, [])
  | `App (`Float, []), `App (`Float, [])
  | `App (`String, []), `App (`String, [])
  | `App (`Range, []), `App (`Range, []) -> ()

  | `App (`List, [ex]), `App (`List, [ac])
  | `App (`Option, [ex]), `App (`Option, [ac]) -> unify ~ex ~ac

  | `App (`Tuple, exs), `App (`Tuple, acs)
  | `App (`Fun, exs), `App (`Fun, acs)
    when List.length exs = List.length acs ->
    List.iter2_exn exs acs ~f:(fun ex ac -> unify ~ex ~ac)

  | `Meta ex, `Meta ac when phys_equal ex ac ->
    ()
                                        (*
  | `Var({ contents = Some(t1') }), _ -> unify t1' t2
  | _, `Var({ contents = Some(t2') }) -> unify t1 t2'
  | `Var({ contents = None } as r1), _ -> (* ������̤����η��ѿ��ξ�� (caml2html: typing_undef) *)
    if occur r1 t2 then raise (Unify_error(t1, t2));
    r1 := Some(t2)
  | _, `Var({ contents = None } as r2) ->
    if occur r2 t1 then raise (Unify_error(t1, t2));
    r2 := Some(t1)
                                         *)
  | _, _ ->
    raise (Unify_error { uniexn_ex = ex; uniexn_ac = ac })

(* �������롼���� (caml2html: typing_g) *)
let rec infer env (e:Ast.t) : (Type.env * Type.t) =
  Printf.printf "infer e: ";
  Ast.print e;
  let unit = Type.(create e.loc desc_unit) in
  let easy_infer env e = snd @@ infer env e in
  let infer_block env es =
    List.fold_left es ~init:(env, unit)
      ~f:(fun (env, _) e -> infer env e)
  in
  try
    let env, desc = match e.desc with
      | `Nop -> (env, desc_unit)

      | `Chunk es -> 
        let env, ty = List.fold_left es
            ~init:(env, unit)
            ~f:(fun (env, ty) e -> infer env e)
        in
        (env, ty.desc)

      | `Return e -> (env, (easy_infer env e).desc)

      | `Unit -> (env, desc_unit)
      | `Bool _ -> (env, desc_bool)
      | `Int _ -> (env, desc_int)
      | `Float _ -> (env, desc_float)
      | `String _ -> (env, desc_string)
      | `Range _ -> (env, desc_range)

      | `List es ->
        begin match es with
          | [] -> (env, desc_list @@ create_metavar e.loc)
          | e :: es ->
            let base_ty = easy_infer env e in
            List.iter es ~f:(fun e ->
                unify ~ex:base_ty ~ac:(easy_infer env e));
            (env, desc_list @@ base_ty)
        end

      | `Tuple es ->
        (env, desc_tuple (List.map es ~f:(easy_infer env)))

      | `Fundef fdef ->
        let params = List.map fdef.fdef_params
            ~f:(fun param -> Type.create_metavar param.loc)
        in
        let ret = Type.create_metavar e.loc in
        let desc = desc_fun params ret in
        let ty = Type.create e.loc desc in
        let env = Env.add_attr env fdef.fdef_name.desc ty in
        let fenv = List.fold2_exn fdef.fdef_params params ~init:env
            ~f:(fun env name ty -> Env.add_attr env name.desc ty)
        in
        let _, ret' = infer_block fenv fdef.fdef_block in
        begin match ret.desc with
          | `Var { contents = Some ret } -> unify ~ex:ret ~ac:ret'
          | `Var ({ contents = None } as ref) -> ref := Some ret'
          | _ -> failwith "return type must be type variable"
        end;
        (env, desc)

      | `Funcall call ->
        let ex_fun = easy_infer env call.fc_fun in
        let args = List.map call.fc_args ~f:(easy_infer env) in
        let ret = Type.create_metavar e.loc in
        let ac_fun = Type.create e.loc (desc_fun args ret) in
        unify ~ex:ex_fun ~ac:ac_fun;
        (env, ex_fun.desc)

      | `Var path ->
        begin match Ast.(path.np_prefix) with
          | Some _ -> failwith "not yet supported"
          | None -> (env, path.np_type.desc) (* TODO: ���Ķ����鸡������ *)
        end
                (*
    | Not(e) ->
      unify `Bool (infer env e);
      `Bool
    | Neg(e) ->
      unify `Int (infer env e);
      `Int
    | Add(e1, e2) | Sub(e1, e2) -> (* ­�����ʤȰ������ˤη����� (caml2html: typing_add) *)
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
    | Let((x, t), e1, e2) -> (* let�η����� (caml2html: typing_let) *)
      unify t (infer env e1);
      g (M.add x t env) e2
    | Var(x) when M.mem x env -> M.find x env (* �ѿ��η����� (caml2html: typing_var) *)
    | Var(x) when M.mem x !extenv -> M.find x !extenv
    | Var(x) -> (* �����ѿ��η����� (caml2html: typing_extvar) *)
      Format.eprintf "free variable %s assumed as external@." x;
      let t = `gentyp () in
      extenv := M.add x t !extenv;
      t
    | LetRec({ name = (x, t); args = yts; body = e1 }, e2) -> (* let rec�η����� (caml2html: typing_letrec) *)
      let env = M.add x t env in
      unify t (`Fun(List.map snd yts, g (M.add_list yts env) e1));
      infer env e2
    | App(e, es) -> (* �ؿ�Ŭ�Ѥη����� (caml2html: typing_app) *)
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
    let ty = Type.create e.loc desc in
    Printf.printf "inferred type: %s\n" (Type.to_string ty);
    (env, ty)
  with
  | Unify_error { uniexn_ex = ex; uniexn_ac = ac } ->
    raise @@ Type_mismatch {
      mismatch_node = deref_term (Type.Env.create ()) e;
      mismatch_ex = deref_type (Type.Env.create ()) ex;
      mismatch_ac = deref_type (Type.Env.create ()) ac;
    }

let run (e:Ast.t) : Ast.t =
  ignore @@ infer (Type.Env.create ()) e;
  deref_term (Type.Env.create ()) e

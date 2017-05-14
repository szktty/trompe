open Core.Std
open Located
open Type
open Logging

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

let tyvar_names = [|
  "a"; "b"; "c"; "d"; "e"; "f"; "g"; "h"; "i"; "j"; "k"; "l"; "m"; "n";
  "o"; "p"; "q"; "r"; "s"; "t"; "u"; "v"; "w"; "x"; "y"; "z"
|]

let rec generalize (ty:Type.t) : Type.t =
  let tyvars : (metavar * tyvar) list ref = ref [] in

  let new_tyvar () =
    Array.get tyvar_names @@ List.length !tyvars
  in

  let rec walk ty =
    let gen = match ty.desc with
      | `App (tycon, args) ->
        let gen_tycon = match tycon with
          | `Tyfun (tyvars, ty) -> `Tyfun (tyvars, generalize ty)
          | tycon -> tycon
        in
        `App (gen_tycon, List.map args ~f:walk)
      | `Var tyvar -> `Var tyvar
      | `Meta ({ contents = None } as ref) ->
        let tyvar = match List.Assoc.find !tyvars ref ~equal:phys_equal with
          | Some tyvar -> tyvar
          | None ->
            let tyvar = new_tyvar () in
            tyvars := (ref, tyvar) :: !tyvars;
            tyvar
        in
        `Var tyvar
      | `Meta ({ contents = Some ty }) -> (walk ty).desc
      | `Poly (tyvars, ty) -> `Poly (tyvars, walk ty)
    in
    Located.create ty.loc gen
  in

  let gen = walk ty in
  if List.length !tyvars = 0 then
    gen
  else begin
    let tyvars = List.rev_map !tyvars ~f:(fun (_, tyvar) -> tyvar) in
    Located.create ty.loc @@ `Poly (tyvars, gen)
  end

let rec subst (ty:Type.t) (env:(tyvar * t) list) =
  match ty.desc with
  | `Var tyvar ->
    Option.value (List.Assoc.find env tyvar) ~default:ty
  | `App (`Tyfun (tyvars, ty'), args) ->
    let env' = List.map2_exn tyvars args
        ~f:(fun tyvar ty -> (tyvar, ty)) in
    subst (subst ty' env') env
  | `App (tycon, args) ->
    let args' = List.map args
        ~f:(fun ty' -> subst ty' env) in
    Located.create ty.loc @@ `App (tycon, args')
  | `Poly (tyvars, ty') ->
    let tyvars' = List.mapi tyvars
        ~f:(fun i tyvar ->
            Array.get Type.var_names (i + List.length env))
    in
    let ty'' = subst ty' @@ List.map2_exn tyvars tyvars'
        ~f:(fun tyvar tyvar' ->
            (tyvar, Located.create ty.loc @@ `Var tyvar'))
    in
    Located.create ty.loc @@ `Poly (tyvars', subst ty'' env)
  | `Meta { contents = Some ty' } ->
    Option.value_map (List.find env ~f:(fun (_, ty) -> ty = ty'))
      ~f:(fun (_, ty') -> subst ty' env)
      ~default:ty
  | `Meta { contents = None } -> ty

let instantiate (ty:Type.t) =
  match ty.desc with
  | `Poly (tyvars, ty') ->
    let env = List.map tyvars
        ~f:(fun tyvar -> (tyvar, create_metavar ty.loc)) in
    subst ty' env
  | _ -> ty

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

(* 型が合うように、型変数への代入をする (caml2html: typing_unify) *)
let rec unify ~(ex:Type.t) ~(ac:Type.t) : unit =
  Printf.printf "unify %s and %s\n" (Type.to_string ex) (Type.to_string ac);
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

  | `Poly _, _ -> unify ~ex:(instantiate ex) ~ac
  | _, `Poly _ -> unify ~ex ~ac:(instantiate ac)

  | `Meta ex, `Meta ac when phys_equal ex ac -> ()
  | `Meta ({ contents = None } as ref), _ -> ref := Some ac
  | _, `Meta ({ contents = None } as ref) -> ref := Some ex
  | `Meta { contents = Some ex }, _ -> unify ~ex ~ac
  | _, `Meta { contents = Some ac } -> unify ~ex ~ac

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
let rec infer env (e:Ast.t) : (Type.t Env.t * Type.t) =
  Printf.printf "infer e: ";
  Ast.print e;
  let unit = Type.(create e.loc desc_unit) in
  let infer_block env es =
    List.fold_left es ~init:(env, unit)
      ~f:(fun (env, _) e -> infer env e)
  in
  try
    let env, desc = match e.desc with
      | `Nop -> (env, desc_unit)

      | `Chunk es -> 
        let env = List.fold_left es
            ~init:env
            ~f:(fun env e -> fst @@ infer env e)
        in
        (env, unit.desc)

      | `Return e -> (env, (easy_infer env e.exp).desc)

      | `If if_ -> 
        let value = Type.create_metavar e.loc in

        let unify_block block =
          List.iteri block ~f:(fun i exp ->
              let ex =
                if i < List.length block - 1 then Type.unit else value
              in
              unify ~ex ~ac:(easy_infer env exp))
        in

        List.iter if_.if_actions
          ~f:(fun (cond, block) ->
              unify ~ex:Type.bool ~ac:(easy_infer env cond);
              unify_block block);
        unify_block if_.if_else;
        (env, value.desc)

      | `Unit -> (env, desc_unit)
      | `Bool _ -> (env, desc_bool)
      | `Int _ -> (env, desc_int)
      | `Float _ -> (env, desc_float)
      | `String _ -> (env, desc_string)
      | `Range _ -> (env, desc_range)

      | `List es ->
        begin match es.exp_list with
          | [] -> (env, desc_list @@ create_metavar e.loc)
          | e :: es ->
            let base_ty = easy_infer env e in
            List.iter es ~f:(fun e ->
                unify ~ex:base_ty ~ac:(easy_infer env e));
            (env, desc_list @@ base_ty)
        end

      | `Tuple es ->
        (env, desc_tuple (List.map es.exp_list ~f:(fun e -> easy_infer env e)))

      | `Fundef fdef ->
        let params = List.map fdef.fdef_params
            ~f:(fun param -> Type.create_metavar param.loc)
        in
        let ret = Type.create_metavar e.loc in
        let fun_ty = Type.fun_ e.loc params ret in
        (* for recursive call *)
        let env = Env.add env fdef.fdef_name.desc fun_ty in
        let fenv = List.fold2_exn fdef.fdef_params params ~init:env
            ~f:(fun env name ty -> Env.add env name.desc ty)
        in
        let _, ret' = infer_block fenv fdef.fdef_block in
        unify ~ex:ret ~ac:ret';
        (env, (generalize fun_ty).desc)

      | `Funcall call ->
        Printf.printf "# funcall ";
        Ast.print e;
        let ex_fun = easy_infer env call.fc_fun in
        let args = List.map call.fc_args ~f:(fun e -> easy_infer env e) in
        let ret = Type.create_metavar e.loc in
        let ac_fun = Type.create e.loc (desc_fun args ret) in
        Printf.printf "# funcall infer ex: %s\n" (Type.to_string ex_fun);
        unify ~ex:ex_fun ~ac:ac_fun;
        Printf.printf "# end funcall infer\n";
        (env, (Type.fun_return ex_fun).desc)

      | `Var path ->
        begin match Ast.(path.np_prefix) with
          | Some _ -> failwith "not yet supported"
          | None ->
            let name = path.np_name.desc in
            match Env.find env name with
            | None -> failwith ("variable is not found: " ^ name)
            | Some ty -> (env, ty.desc)
        end

      | `Unexp exp ->
        let op_ty, val_ty = match exp.unexp_op.desc with
          | `Pos | `Neg -> (Type.int, Type.int)
          | `Fpos | `Fneg -> (Type.float, Type.float)
          | _ -> failwith "not yet supported"
        in
        unify op_ty (easy_infer env e);
        (env, val_ty.desc)

      | `Binexp { binexp_left = e1; binexp_op = op; binexp_right = e2 } ->
        let op_ty, val_ty = match op.desc with
          | `Eq | `Ne ->
            (easy_infer env e1, Type.bool)
          | `And | `Or ->
            (Type.bool, Type.bool)
          | `Lt | `Le | `Gt | `Ge ->
            (Type.int, Type.bool)
          | `Add | `Sub | `Mul | `Div
          | `Pow | `Mod | `Lcomp | `Rcomp ->
            (Type.int, Type.int)
          | `Fadd | `Fsub | `Fmul | `Fdiv ->
            (Type.float, Type.float)
          | _ -> failwith "not yet supported"
        in
        unify op_ty (easy_infer env e1);
        unify op_ty (easy_infer env e2);
        (env, val_ty.desc)

                (*
    | Not(e) ->
      unify `Bool (infer env e);
      `Bool
    | Neg(e) ->
      unify `Int (infer env e);
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
    let ty = Type.create e.loc desc in
    Printf.printf "inferred node: ";
    Ast.print e;
    Printf.printf "inferred type: %s\n" (Type.to_string ty);
    (env, ty)
  with
  | Unify_error { uniexn_ex = ex; uniexn_ac = ac } ->
    raise @@ Type_mismatch {
      mismatch_node = e;
      mismatch_ex = generalize ex;
      mismatch_ac = generalize ac;
    }

and easy_infer env (e:Ast.t) : Type.t = snd @@ infer env e

let run (e:Ast.t) =
  verbosef "begin typing";
  ignore @@ infer (Runtime.type_env ()) e;
  verbosef "end typing"

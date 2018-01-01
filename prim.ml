open Base
open Runtime

type t = Runtime.prim

let call rt ~ctx ~prim ~args ~tys =
  let arity = List.length args in
  if prim.prim_arity <> arity then
    Error (Invalid_arity (arity, prim.prim_arity))
  else begin
    let args = Args.create args tys in
    prim.prim_fun rt ctx args
  end

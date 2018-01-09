let rec eval rt ~ctx ~node =
  match node with
  | Ast.String tok ->
    rt, ctx, Value.String tok.text
  | _ ->
    rt, ctx, Value.Void

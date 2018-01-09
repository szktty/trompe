val eval : Runtime.t ->
  ctx:Runtime.context ->
  node:Ast.t ->
  (Runtime.t * Runtime.context * Value.t)

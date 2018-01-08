let last_cmt : Token.comment option ref = ref None

let token ?cmt ?next ~pos text =
  { Token.cmt; pos; text; next }

open Base

type t = {
  line : int;
  col : int;
  offset : int;
}

let zero = { line = 0; col = 0; offset = 0 }

let of_lexing_pos (pos : Lexing.position) =
  { line = pos.pos_lnum;
    col = pos.pos_bol;
    offset = pos.pos_cnum;
  }

let equal pos1 pos2 =
  (pos1.line = pos2.line) &&
  (pos1.col = pos2.col) &&
  (pos1.offset = pos2.offset)

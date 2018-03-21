type comment = {
  cmt_pos : Position.t;
  cmt_text : string;
}

type t = {
  pos : Position.t;
  cmt : comment option;
  text : string;
  next : t option;
}

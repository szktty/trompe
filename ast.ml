type t =
  | Nop (* for debug *)
  | String of Token.t

let start_pos = function
  | String tok -> tok.pos
  | _ -> failwith "notimpl"

let end_pos = function
  | _ -> failwith "notimpl"

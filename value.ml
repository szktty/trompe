open Base

type t =
  | Void
  | Int of int
  | String of string
  | Prim of string
  | Clos of closure
  | Fun of t

and closure = {
  clos_env : t Map.M(String).t;
}

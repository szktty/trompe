open Base

type t =
  | Void
  | Bool of bool
  | Int of int
  | String of string
  | List of t list
  | Tuple of t list
  | Prim of string
  | Clos of closure
  | Fun of t

and closure = {
  clos_env : t Map.M(String).t;
}

let type_name = function
  | Int _ -> "int"
  | String _ -> "string"
  | _ -> failwith "notimpl"

let string = function
  | String s -> Some s
  | _ -> None

let string_exn v =
  match v with
  | String s -> s
  | _ -> raise (Runtime_exc.invalid_type ~given:(type_name v) ~valid:"string")

let list = function
  | List es -> Some es
  | _ -> None

let list_exn v =
  match v with
  | List es -> es
  | _ -> raise (Runtime_exc.invalid_type ~given:(type_name v) ~valid:"list")


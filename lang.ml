open Core.Std

type t = [
  | `Module of t Module.t
  | `Unit
  | `Bool of bool
  | `String of string
  | `Int of int
  | `Float of float
  | `Range of (int * int)
  | `List of t list
  | `Tuple of t list
  | `Fun of (Ast.fundef * t Env.t)
  | `Prim of string
  | `Ref of t ref
  | `Enum of (string, t) Ast.enum
  | `Exn of user_error
]

and primitive = context -> t Env.t -> t list -> t

and context = {
  ctx_parent : context option;
  ctx_call : Ast.t option;
}

and l_exn = {
  exn_desc : exn_desc;
  exn_reason : string option;
}

and exn_desc =
  | Fatal_error
  | Name_error
  | Runtime_error
  | Standard_error
  | Value_error
  | User_error of user_error

and user_error = {
  user_error_name : string;
  user_error_value : t option;
  user_error_reason : string option;
}

let rec to_string value =

  let open Printf in
  let values_to_string values open_tag close_tag sep =
    let buf = Buffer.create 16 in
    Buffer.add_string buf open_tag;
    List.iter values ~f:(fun v ->
        Buffer.add_string buf (to_string v);
        Buffer.add_string buf sep;
        Buffer.add_string buf " ");
    Buffer.add_string buf close_tag;
    Buffer.contents buf

  in
  match value with
  | `Unit -> "()"
  | `Bool true -> "true"
  | `Bool false -> "false"
  | `String v -> sprintf "\"%s\"" v
  | `Int v -> Int.to_string v
  | `Float v -> Float.to_string v
  | `List vs -> values_to_string vs "[" "]" ";"
  | `Tuple vs -> values_to_string vs "(" ")" ","
  | `Fun (fdef, _) -> sprintf "fun/%d" (List.length Ast.(fdef.fdef_params))
  | `Prim name -> sprintf "#primitive(\"%s\")" name
  | _ -> "?"

module Context = struct

  type t = context

  let create parent call =
    { ctx_parent = parent; ctx_call = call }

end

module Exn = struct

  type t = l_exn

  let name e =
    match e.exn_desc with
    | Fatal_error -> "Fatal_error"
    | Name_error -> "Name_error"
    | Runtime_error -> "Runtime_error"
    | Standard_error -> "Standard_error"
    | Value_error -> "Value_error"
    | User_error e -> e.user_error_name

  let reason e =
    match e.exn_desc with
    | User_error e -> e.user_error_reason
    | _ -> e.exn_reason

  let of_user_error e =
    { exn_desc = User_error e; exn_reason = e.user_error_reason }

  let of_reason e reason =
    { exn_desc = e; exn_reason = Some reason }

end

module Op = struct

  type compare =
    | Asc
    | Desc
    | Equal

  let compare x y =
    match (x, y) with
    | (`Int x', `Int y') ->
      if x' = y' then
        Equal
      else if x' < y' then
        Asc
      else
        Desc
    | _ -> failwith "not supported value types"

  let le x y =
    match compare x y with
    | Equal | Asc -> `Bool true
    | Desc -> `Bool false

  let compute x y ~f_int ~f_float =
    match (x, y) with
    | (`Int x', `Int y') -> `Int (f_int x' y')
    | (`Float x', `Float y') -> `Float (f_float x' y')
    | _ -> failwith "must integer or float values"

  let add x y = compute x y ~f_int:(+) ~f_float:(+.)
  let sub x y = compute x y ~f_int:(-) ~f_float:(-.)
  let mul x y = compute x y ~f_int:( * ) ~f_float:( *. )
  let div x y = compute x y ~f_int:(/) ~f_float:(/.)
  let mod_ x y = compute x y ~f_int:(mod)
      ~f_float:(fun _ _  -> failwith "values must integer for % operator")

end

let rec equal x y =
  match (x, y) with
  | (`Unit, `Unit) -> true
  | (`Bool x, `Bool y) -> x = y
  | (`String x, `String y) -> x = y
  | (`Int x, `Int y) -> x = y
  | (`Float x, `Float y) -> x =. y
                                   (*
  | L_range of (int * int)
  | L_list of t list
                                    *)
  | (`Tuple xs, `Tuple ys) ->
    if List.length xs <> List.length ys then
      false
    else
      List.for_all2_exn xs ys ~f:equal
                                   (*
  | L_fun of (Ast.fundef * env)
  | L_prim of (string * primitive)
                                    *)
  | _ -> failwith "Lang.equal not support"

open Core.Std

type t = [
  | `Module of string
  | `Unit
  | `Bool of bool
  | `String of string
  | `Int of int
  | `Float of float
  | `Range of (int * int)
  | `List of t list
  | `Tuple of t list
  | `Fun of (Ast.fundef * capture)
  | `Stream of (In_channel.t option * Out_channel.t option)
  | `Prim of string
  | `Ref of t ref
  | `Enum of (string, t) Ast.enum
  | `Exn of user_error
]

and value = t

and capture = t String.Map.t

and user_error = {
  user_error_name : string;
  user_error_value : t option;
  user_error_reason : string option;
}

and primitive = t list -> t

let rec to_string value =
  let open Printf in
  let values_to_string values open_tag close_tag sep =
    let buf = Buffer.create 16 in
    Buffer.add_string buf open_tag;
    String.concat (List.map values ~f:to_string) ~sep:", "
    |> Buffer.add_string buf;
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
  | `List vs -> values_to_string vs "[" "]" ","
  | `Tuple vs -> values_to_string vs "(" ")" ","
  | `Fun (fdef, _) -> sprintf "fun/%d" (List.length Ast.(fdef.fdef_params))
  | `Prim name -> sprintf "#primitive(\"%s\")" name
  | _ -> "?"

module Context = struct

  type t = {
    belong : t Module.t option;
    parent : t option;
    callee : Ast.t option;
  }

  let create ?(belong=None) ?(parent=None) ?(callee=None) () =
    { belong; parent; callee }

end

type l_exn = {
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

  let to_bool_value value = `Bool value

  let eq x y =
    to_bool_value (x = y)

  let ne x y =
    to_bool_value (x <> y)

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
  | _ -> failwith "Value.equal not support"

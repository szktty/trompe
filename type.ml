open Base

type name = string

type t = desc Located.t

and desc =
  | App of tycon * t list
  | Var of name
  | Poly of name list * t
  | Meta of t option ref

and tycon =
  | Unit
  | Bool
  | Int
  | Float
  | String
  | List
  | Tuple
  | Range
  | Option
  | Box
  | Struct of name list
  | Enum of name list
  | Fun
  | Fun_printf
  | Module of name
  | Stream
  | Tyfun of name list * t
  | Unique of name

let tyvar_names = [|
  "a"; "b"; "c"; "d"; "e"; "f"; "g"; "h"; "i"; "j"; "k"; "l"; "m"; "n";
  "o"; "p"; "q"; "r"; "s"; "t"; "u"; "v"; "w"; "x"; "y"; "z";
|]

let tyvar_name n =
  Array.get tyvar_names n

let create loc desc =
  Located.create desc ~loc

let meta loc =
  Located.create (Meta (ref None)) ~loc

let equal (ty1:t) (ty2:t) =
  Caml.(ty1.desc = ty2.desc)

let rec to_string (ty:t) =
  match ty.desc with
  | App (tycon, args) ->
    let tycon_s = match tycon with
      | Unit -> "Unit"
      | Bool -> "Bool"
      | Int -> "Int"
      | Float -> "Float"
      | String -> "String"
      | List -> "List"
      | Tuple -> "Tuple"
      | Range -> "Range"
      | Fun -> "Fun"
      | Fun_printf -> "Fun_printf"
      | Stream -> "Stream"
      | Option -> "Option"
      | Box -> "Box"
      | Module name -> Printf.sprintf "Module(%s)" name
      | _ -> failwith "not impl"
    in
    let args_s =
      List.map args ~f:to_string
      |> String.concat ~sep:", "
    in
    Printf.sprintf "App(%s, [%s])" tycon_s args_s
  | Meta { contents = None } -> "Meta(?)"
  | Meta { contents = Some ty } -> Printf.sprintf "Meta(%s)" (to_string ty)
  | Var name -> Printf.sprintf "Var(%s)" (Utils.quote name)
  | Poly (tyvars, ty) ->
    Printf.sprintf "Poly([%s], %s)" 
      ((String.concat (List.map tyvars ~f:Utils.quote) ~sep:", "))
      (to_string ty)

module Spec = struct

  let tyvar name = Located.create @@ Var (tyvar_name name)

  let tyvar_a = tyvar 0

  let tyvar_b = tyvar 1

  let tyvar_c = tyvar 2

  let tyvar_d = tyvar 3

  let poly tyvars ty = Poly (tyvars, ty)

  let app ?(args=[]) tycon = App (tycon, args)

  let desc_unit = app Unit

  let desc_bool = app Bool

  let desc_int = app Int

  let desc_float = app Float

  let desc_string = app String

  let desc_range = app Range

  let desc_list e = app ~args:[e] List

  let desc_tuple es = app ~args:es Tuple

  let desc_option e = app ~args:[e] Option

  let desc_box e = app ~args:[e] Box

  let desc_fun params ret =
    app ~args:(List.append params [ret]) Fun

  let desc_fun_printf = app Fun_printf

  let desc_stream = app Stream

  let unit = Located.create desc_unit

  let bool = Located.create desc_bool

  let int = Located.create desc_int

  let float = Located.create desc_float

  let string = Located.create desc_string

  let range = Located.create desc_range

  let list e = Located.create @@ desc_list e

  let list_gen = Located.create @@ poly ["a"] (list tyvar_a)

  let tuple es = Located.create @@ desc_tuple es

  let option e = Located.create @@ desc_option e

  let option_gen = Located.create @@ poly ["a"] (option tyvar_a)

  let box e = Located.create @@ desc_box e

  let box_gen = Located.create @@ poly ["a"] (box tyvar_a)

  let fun_ loc params ret = Located.create ~loc @@ desc_fun params ret

  let fun_printf = Located.create @@ desc_fun_printf

  let module_ name = Located.create @@ app (Module name)

  let stream = Located.create @@ desc_stream 

  let unique loc ty name =
    Located.create ~loc (app (Unique name) ~args:[ty])

  let struct_ loc fields =
    let names, tys = List.fold_left fields ~init:([], [])
        ~f:(fun (names, tys) (name, ty) ->
            name :: names, ty :: tys)
    in
    Located.create ~loc (app (Struct (List.rev names)) ~args:(List.rev tys))

end

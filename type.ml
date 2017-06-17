open Core.Std

type t = desc Located.t

and desc = [
  | `App of tycon * t list
  | `Var of tyvar
  | `Poly of tyvar list * t
  | `Meta of metavar
]

and tycon = [
  | `Unit
  | `Bool
  | `Int
  | `Float
  | `String
  | `List
  | `Tuple
  | `Range
  | `Option
  | `Struct of string list
  | `Enum of string list
  | `Fun
  | `Fun_printf
  | `Stream
  | `Tyfun of tyvar list * t
  | `Unique of tycon * int
]

and tyvar = string

and metavar = t option ref

let create loc desc =
  Located.create loc desc

let create_metavar loc =
  Located.create loc (`Meta (ref None))

let var_names = [|
  "a"; "b"; "c"; "d"; "e"; "f"; "g"; "h"; "i"; "j"; "k"; "l"; "m"; "n";
  "o"; "p"; "q"; "r"; "s"; "t"; "u"; "v"; "w"; "x"; "y"; "z";
|]

let equal (ty1:t) (ty2:t) =
  ty1.desc = ty2.desc

let rec to_string (ty:t) =
  match ty.desc with
  | `App (tycon, args) ->
    let tycon_s = match tycon with
      | `Unit -> "Unit"
      | `Bool -> "Bool"
      | `Int -> "Int"
      | `Float -> "Float"
      | `String -> "String"
      | `List -> "List"
      | `Tuple -> "Tuple"
      | `Range -> "Range"
      | `Fun -> "Fun"
      | `Fun_printf -> "Fun_printf"
      | `Stream -> "Stream"
      | `Option -> "Option"
      | _ -> failwith "not impl"
    in
    let args_s =
      List.map args ~f:to_string
      |> String.concat ~sep:", "
    in
    Printf.sprintf "App(%s, [%s])" tycon_s args_s
  | `Meta { contents = None } -> "Meta(?)"
  | `Meta { contents = Some ty } -> "Meta(" ^ to_string ty ^ ")"
  | `Var name -> "Var(" ^ name ^ ")"
  | `Poly (tyvars, ty) ->
    "Poly([" ^ (String.concat tyvars ~sep:", ") ^ "], " ^ to_string ty ^ ")"

let rec fun_return (ty:t) =
  match ty.desc with
  | `Meta { contents = Some ty }
  | `Poly (_, ty) -> fun_return ty
  | `App (`Fun, args) -> List.last_exn args
  | _ -> failwith "not function"

let app ?(args=[]) tycon = `App (tycon, args)
let desc_unit = app `Unit
let desc_bool = app `Bool
let desc_int = app `Int
let desc_float = app `Float
let desc_string = app `String
let desc_range = app `Range
let desc_list e = app ~args:[e] `List
let desc_tuple es = app ~args:es `Tuple
let desc_option e = app ~args:[e] `Option
let desc_fun params ret =
  app ~args:(List.append params [ret]) `Fun
let desc_fun_printf = app `Fun_printf
let desc_stream = app `Stream

let unit = Located.less desc_unit
let bool = Located.less desc_bool
let int = Located.less desc_int
let float = Located.less desc_float
let string = Located.less desc_string
let range = Located.less desc_range
let list e = Located.less @@ desc_list e
let tuple es = Located.less @@ desc_tuple es
let option e = Located.less @@ desc_option e
let fun_ loc params ret = Located.create loc @@ desc_fun params ret
let fun_printf = Located.less @@ desc_fun_printf
let stream = Located.less @@ desc_stream 

let parse_format s =
  let module F = Utils.Format in
  let fmt = Utils.Format.parse s in
  List.map (F.params fmt)
    ~f:(function
        | F.Int -> int
        | F.String -> string
        | _ -> failwith "parse_format")

module Spec = struct

  type t = [
    | `Tyvar of string
    | `Unit
    | `Bool
    | `Int
    | `Float
    | `String
    | `List of t
    | `Tuple of t list
    | `Range
    | `Option of t
    | `Fun of t list
    | `Fun_printf
    | `Stream
  ]

  let unit = `Unit
  let bool = `Bool
  let int = `Int
  let float = `Float
  let string = `String
  let list e = `List e
  let tuple es = `Tuple es
  let range = `Range
  let option e = `Option e
  let fun_printf = `Fun_printf
  let stream = `Stream

  let a = `Tyvar "a"
  let b = `Tyvar "b"
  let c = `Tyvar "c"
  let d = `Tyvar "d"
  let e = `Tyvar "e"

  let (@->) x y =
    match x with
    | `Fun args -> `Fun (List.append args [y])
    | _ -> `Fun [x; y]

  let flat_tyvars tyvars =
    let tyvars =
      List.fold_left tyvars
        ~init:[]
        ~f:(fun accu tyvar ->
            if List.existsi accu ~f:(fun _ e -> e = tyvar) then
              accu
            else
              tyvar :: accu)
    in
    List.sort tyvars ~cmp:String.Caseless.descending

  let collect_tyvars (spec:t) =
    let rec f (tyvars:string list) spec =
      match spec with
      | `Unit -> tyvars, desc_unit
      | `Bool -> tyvars, desc_bool
      | `Int -> tyvars, desc_int
      | `Float -> tyvars, desc_float
      | `String -> tyvars, desc_string
      | `Stream -> tyvars, desc_stream
      | `Fun_printf -> tyvars, desc_fun_printf
      | `Tyvar name -> (name :: tyvars), `Var name
      | `List e ->
        let tyvars', ty = f tyvars e in
        tyvars', desc_list (Located.less ty)
      | `Fun args ->
        let tyvars', args' =
          List.fold_left args ~init:(tyvars, [])
            ~f:(fun (tyvars, args) arg ->
                let tyvars', arg' = f tyvars arg in
                tyvars', Located.less arg' :: args)
        in
        tyvars', `App (`Fun, List.rev args')
      | _ -> failwith "not yet support"
    in
    f [] spec

  let to_type spec =
    match collect_tyvars spec with
    | [], desc -> Located.less desc
    | tyvars, desc ->
      let tyvars = flat_tyvars tyvars in
      Located.less @@ `Poly (tyvars, Located.less desc)

end

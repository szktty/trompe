open Core.Std

type t = desc Located.t

and _t = t

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
  | `Tyfun of tyvar list * t
  | `Unique of tycon * int
]

and tyvar = string

and metavar = t option ref

module Env = Env.Make(struct
    type t = _t
  end)

module Module = Module.Make(Env)


let create loc desc =
  Located.create loc desc

let create_metavar loc =
  Located.create loc (`Meta (ref None))

let var_names = [|
  "a"; "b"; "c"; "d"; "e"; "f"; "g"; "h"; "i"; "j"; "k"; "l"; "m"; "n";
  "o"; "p"; "q"; "r"; "s"; "t"; "u"; "v"; "w"; "x"; "y"; "z";
|]

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
      | `Option -> "Option"
      | _ -> failwith "not impl"
    in
    "App(" ^ tycon_s ^ ")"
  | `Meta { contents = None } -> "Meta(_)"
  | `Meta { contents = Some ty } -> "Meta(" ^ to_string ty ^ ")"
  | `Var name -> "Var(" ^ name ^ ")"
  | `Poly (tyvars, ty) ->
    "Poly([" ^ (String.concat tyvars ~sep:", ") ^ "], " ^ to_string ty ^ ")"

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

  let a = `Tyvar "a"
  let b = `Tyvar "b"
  let c = `Tyvar "c"
  let d = `Tyvar "d"
  let e = `Tyvar "e"

  let (+>) x y =
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

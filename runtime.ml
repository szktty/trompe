open Core.Std

let type_modules : Type.t Module.t list ref = ref []

let value_modules : Value.t Module.t list ref = ref []

let find_module ?(path=[]) tops ~name =
  match path with
  | [] -> List.find tops ~f:(fun m -> Module.name m = name)
  | fst :: rest ->
    match List.find tops ~f:(fun m -> Module.name m = name) with
    | None -> None
    | Some m -> Module.find_module m ~prefix:rest ~name

let find_type_module path =
  find_module !type_modules path

let find_value_module path =
  find_module !value_modules path

module Primitive = struct

  let primitives : Value.primitive String.Map.t ref = ref String.Map.empty

  let add name f =
    primitives := String.Map.add !primitives ~key:name ~data:f

  let find name =
    String.Map.find !primitives name

end

module Spec = struct

  type attr = {
    attr_name : string;
    ty : Type.t;
    value : Value.t;
  }

  type t = {
    mod_name : string;
    parent : string option;
    attrs : attr list;
  }

  let define ?parent name =
    { mod_name = name; parent; attrs = [] }

  let (+>) def attr =
    { def with attrs = attr :: def.attrs }

  let attr name ty value = 
    { attr_name = name; ty; value }

  let int name value =
    attr name Type.int (`Int value)

  let string name value =
    attr name Type.string (`String value)

  let fun_ name (spec:Type.Spec.t) pname =
    attr name (Type.Spec.to_type spec) (`Prim pname)

  let end_ spec =
    (* TODO: parent *)
    let tattrs, vattrs = List.fold_left spec.attrs
        ~init:(String.Map.empty, String.Map.empty)
        ~f:(fun (tattrs, vattrs) attr ->
            (String.Map.add tattrs ~key:attr.attr_name ~data:attr.ty,
             String.Map.add vattrs ~key:attr.attr_name ~data:attr.value))
    in
    type_modules := Module.create spec.mod_name ~attrs:tattrs :: !type_modules;
    value_modules := Module.create spec.mod_name ~attrs:vattrs :: !value_modules;
    ()

end

(*
let test () =
  let kernel = Spec.(define "Kernel"
                     +> fun_ "show" Type.Spec.(a @-> unit) "show"
                     +> string "version" "0.0.1"
                    )
  in
  let sub = Module.define "Test" ~parent:kernel in
  sub
 *)

    (*
let test () =
  Fun.define "show" attr_show
     *)

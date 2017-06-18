open Core.Std

let type_modules : Type.t Module.t list ref = ref []

let type_imports : Type.t Module.t list ref = ref []

let value_modules : Value.t Module.t list ref = ref []

let value_imports : Value.t Module.t list ref = ref []

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

let top_module_names mods =
  List.fold_left mods
    ~init:[]
    ~f:(fun accu m ->
        match Module.parent m with
        | Some _ -> accu
        | None -> Module.name m :: accu)

let top_module_attrs mods ~f =
  List.map (top_module_names mods) ~f:(fun name -> (name, f name))
  |> String.Map.of_alist_exn

let type_env () =
  let attrs = top_module_attrs !type_modules
      ~f:(fun name -> Type.module_ name) in
  Env.create ~imports:!type_modules ~attrs ()

let value_env () =
  let attrs = top_module_attrs !value_modules ~f:(fun name -> `Module name) in
  Env.create ~imports:!value_imports ~attrs ()

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
    init : bool;
    parent : string option;
    attrs : attr list;
  }

  let define ?parent ?(init=false) name =
    { mod_name = name; init; parent; attrs = [] }

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
    Printf.printf "# add module %s\n" spec.mod_name;
    let tattrs, vattrs = List.fold_left spec.attrs
        ~init:(String.Map.empty, String.Map.empty)
        ~f:(fun (tattrs, vattrs) attr ->
            (String.Map.add tattrs ~key:attr.attr_name ~data:attr.ty,
             String.Map.add vattrs ~key:attr.attr_name ~data:attr.value))
    in
    let tmod = Module.create spec.mod_name ~attrs:tattrs in
    let vmod = Module.create spec.mod_name ~attrs:vattrs in
    type_modules := tmod :: !type_modules;
    value_modules := vmod :: !value_modules;
    if spec.init then begin
      type_imports := tmod :: !type_imports;
      value_imports := vmod :: !value_imports;
    end;
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

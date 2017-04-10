open Core.Std

module type S = sig

  type t
  type primitive

  module Env : Env.S

  val toplevel : t list ref

  val define : t -> unit

  val create :
    ?parent:t option
    -> ?submodules:t list
    -> ?imports: t list
    -> ?env:Env.t option
    -> name:string
    -> unit
    -> t

  val name : t -> string

  val root : t -> t option

  val is_root : t -> bool

  val import : t -> t -> unit

  val namepath : t -> Namepath.t

  val find_module :
    ?prefix:string list
    -> t
    -> name:string
    -> (t, t * string) Result.t

  val add_module : t -> t -> unit

  val find_attr : t -> string -> Env.data option

  val add_attr : t -> key:string -> data:Env.data -> unit

  val primitives : unit -> primitive String.Map.t

  val add_primitive : name:string -> primitive:primitive -> unit

  val find_primitive : string -> primitive option

end

module Make(A: sig
    module Env : Env.S
    type primitive
  end) : S with type primitive = A.primitive = struct

  type primitive = A.primitive

  type t = {
    parent : t option;
    name : string;
    mutable env : A.Env.t;
    mutable submodules : t list;
    mutable imports : t list;
  }

  module Env = A.Env

  let toplevel = ref []

  let define m =
    toplevel := m :: !toplevel

  let create ?(parent=None) ?(submodules=[]) ?(imports=[]) ?(env=None) ~name () =
    let env = match env with
      | Some env -> env
      | None -> A.Env.create ()
    in
    { parent; name; env; submodules; imports }

  let name m = m.name

  let rec root m =
    match m.parent with
    | None -> Some m
    | Some m -> root m

  let is_root m = Option.is_none m.parent

  let import m x =
    m.imports <- x :: m.imports

  let rec namepath m =
    match m.parent with
    | None -> Namepath.create m.name
    | Some m ->
      Namepath.create ~prefix:(Some (namepath m)) m.name

  let rec find_module ?(prefix=[]) m ~name =
    match prefix with
    | fst :: rest ->
      begin match find_module m ~name:fst with
        | Result.Error _ as e -> e
        | Result.Ok sub -> find_module ~prefix:rest sub ~name
      end
    | [] ->
      let from = List.append m.submodules !toplevel in
      match List.find from ~f:(fun m -> m.name = name) with
      | None -> Result.Error (m, name)
      | Some m -> Result.Ok m

  let add_module m x =
    m.submodules <- x :: m.submodules

  let rec find_attr m key =
    match A.Env.find m.env key with
    | Some _ as res -> res
    | None ->
      match List.find_mapi m.imports
              ~f:(fun _ m -> find_attr m key) with
      | Some _ as v -> v
      | None -> None

  let add_attr m ~key ~data =
    m.env <- A.Env.add m.env ~key ~data

  let prim_ref : primitive String.Map.t ref = ref String.Map.empty

  let primitives () = !prim_ref

  let add_primitive ~name ~primitive =
    prim_ref := String.Map.add !prim_ref ~key:name ~data:primitive

  let find_primitive name =
    String.Map.find !prim_ref name

end

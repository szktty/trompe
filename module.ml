open Core.Std

module type S = sig

  type t
  type env
  type data

  val toplevel : t list ref

  val define : t -> unit

  val create :
    ?parent:t option
    -> ?submodules:t list
    -> ?imports: t list
    -> ?env:env option
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

  val find_attr : t -> string -> data option

  val add_attr : t -> key:string -> data:data -> unit

end

module Make(E: Env.S) : S = struct

  type env = E.t

  type data = E.data

  type t = {
    parent : t option;
    name : string;
    mutable env : env;
    mutable submodules : t list;
    mutable imports : t list;
  }

  let toplevel = ref []

  let define m =
    toplevel := m :: !toplevel

  let create ?(parent=None) ?(submodules=[]) ?(imports=[]) ?(env=None) ~name () =
    let env = match env with
      | Some env -> env
      | None -> E.create ()
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
    match E.find m.env key with
    | Some _ as res -> res
    | None ->
      match List.find_mapi m.imports
              ~f:(fun _ m -> find_attr m key) with
      | Some _ as v -> v
      | None -> None

  let add_attr m ~key ~data =
    m.env <- E.add m.env ~key ~data

end

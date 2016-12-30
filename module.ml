open Core.Std

module type S = sig

  type t
  type env
  type data

  val create :
    ?parent:t option
    -> ?subs:t list
    -> ?imports: t list
    -> name:string
    -> env:env
    -> t

  val find_sub :
    ?prefix:string list
    -> t
    -> name:string
    -> (t, t * string) Result.t

  val find_attr : t -> string -> data option

  (* TODO:
   * val root : t -> t option
   * val is_root : t -> bool
   * val namepath : t -> Namepath.t
   * val import : t -> t -> unit
   * val add_attr : t -> key:string -> data:data -> unit
  *)
end

module Make(E: Env.S) : S = struct

  type env = E.t

  type data = E.data

  type t = {
    parent : t option;
    name : string;
    mutable env : env;
    mutable subs : t list;
    mutable imports : t list;
  }

                (*
  let top_modules : t list ref = ref []

  let register m =
    top_modules := m :: !top_modules
                 *)

  let create ?(parent=None) ?(subs=[]) ?(imports=[]) ~name ~env =
    { parent; name; env; subs; imports }

  let rec find_sub ?(prefix=[]) m ~name =
    match prefix with
    | fst :: rest ->
      begin match find_sub m ~name:fst with
        | Result.Error _ as e -> e
        | Result.Ok sub -> find_sub ~prefix:rest sub ~name
      end
    | [] ->
      match List.find m.subs ~f:(fun m -> m.name = name) with
      | None -> Result.Error (m, name)
      | Some sub -> Result.Ok sub

  let rec find_attr m key =
    match E.find m.env key with
    | Some _ as res -> res
    | None ->
      match List.find_mapi m.imports
              ~f:(fun _ m -> find_attr m key) with
      | Some _ as v -> v
      | None -> None

end

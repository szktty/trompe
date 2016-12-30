open Core.Std

module type S = sig

  type t

  type data

  val create :
    ?parent:t option
    -> ?attrs:(string * data) list
    -> unit
    -> t

  val find : t -> string -> data option
  val add : t -> key:string -> data:data -> t

  (* TODO:
   * val debug : t -> unit
   * val to_map : t -> data String.Map.t
  *)

end

module Make(A : sig
    type t
    (* TODO: val string_of_data : t -> string *)
  end) : S = struct

  type data = A.t

  type t = {
    parent : t option;
    attrs : A.t String.Map.t;
  }

  let create ?(parent=None) ?(attrs=[]) () =
    { parent = parent;
      attrs = String.Map.of_alist_reduce attrs ~f:(fun _ b -> b);
    }

  let rec find env key =
    match String.Map.find env.attrs key with
    | Some _ as res -> res
    | None ->
      match env.parent with
      | None -> None
      | Some env -> find env key

  let add env ~key ~data =
    { env with attrs = String.Map.add env.attrs ~key ~data }

end

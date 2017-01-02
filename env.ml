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

  val merge : t -> data String.Map.t -> t

  val concat : t -> data String.Map.t

  val debug : t -> f:(data -> string) -> unit

end

module Make(A : sig
    type t
  end) : S with type data = A.t = struct

  type t = {
    parent : t option;
    attrs : A.t String.Map.t;
  }

  type data = A.t

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

  let merge env map =
    { env with attrs = String.Map.merge env.attrs map
                   ~f:(fun ~key owner ->
                       match owner with
                       | `Left v | `Right v -> Some v
                       | `Both (_, v2) -> Some v2)
    }

  let concat env =
    let rec f attrs accu =
      String.Map.fold attrs ~init:accu
        ~f:(fun ~key ~data accu -> String.Map.add accu ~key ~data)
    in
    f env.attrs String.Map.empty

  let debug env ~f =
    let print env indent =
      let open Printf in
      let indent_s = String.make (indent * 2) ' ' in
      printf "%s{\n" indent_s;
      String.Map.iteri env.attrs ~f:(fun ~key ~data ->
          printf "%s  %s = %s\n" indent_s key (f data));
      printf "%s}\n" indent_s
    in
    let rec f env indent =
      let indent = match env.parent with
        | None -> indent
        | Some parent ->
          f env indent;
          indent + 1
      in
      print env indent
    in
    f env 0

end

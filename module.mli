open Core.Std

type 'a t

val create :
  ?parent:'a t
  -> ?submodules:'a t list
  -> ?imports:'a t list
  -> ?attrs:'a String.Map.t
  -> string
  -> 'a t

val name : 'a t -> string

val root : 'a t -> 'a t option

val is_root : 'a t -> bool

val parent : 'a t -> 'a t option

val import : 'a t -> 'a t -> unit

val namepath : 'a t -> string Namepath.t

val find_module : ?prefix:string list -> 'a t -> name:string -> 'a t option

val add_module : 'a t -> 'a t -> unit

val find_attr : 'a t -> string -> 'a option

val add_attr : 'a t -> key:string -> data:'a -> unit


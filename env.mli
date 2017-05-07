open Core.Std

type 'a t

val create :
  ?parent:'a t option
  -> ?attrs:(string * 'a) list
  -> unit
  ->'a t

val find : 'a t -> string -> 'a option

val add : 'a t -> key:string -> data:'a ->'a t

val merge : 'a t -> 'a String.Map.t ->'a t

val concat : 'a t -> 'a String.Map.t

val debug : 'a t -> f:('a -> string) -> unit

open Core.Std

let debugf f =
  if !Config.debug_mode then
    printf ("# " ^^ f ^^ "\n")
  else
    Printf.ifprintf stderr f

let verbosef f =
  if !Config.verbose_mode || !Config.debug_mode then
    printf ("# " ^^ f ^^ "\n")
  else
    Printf.ifprintf stderr f

open Core.Std

let init () =
  Lib_kernel.init ();
  Lib_int.init ();
  ()

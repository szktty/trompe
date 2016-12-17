open Core.Std

let install () =
  Lib_kernel.install ();
  Lib_int.install ();
  ()

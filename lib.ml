open Core.Std

let init () =
  Lib_filename.init ();
  Lib_int.init ();
  Lib_io.init ();
  Lib_kernel.init ();
  Lib_list.init ();
  Lib_os.init ();
  Lib_string.init ();
  ()

open Core.Std
open Value
open Interp.Primitive

let prim_open args : Value.t =
  check_arity "io.open" 1 args;
  let name = get_string args 0 in
  let in_chan = In_channel.create name in
  `Stream (Some in_chan, None)

let prim_close args =
  check_arity "io.close" 1 args;
  let in_, out = get_stream args 0 in
  Option.iter in_ ~f:In_channel.close;
  Option.iter out ~f:In_channel.close;
  `Unit

let prim_read_all args =
  check_arity "io.read_all" 1 args;
  match get_stream args 0 with
  | None, _ -> failwith "not input stream"
  | Some in_, _ -> `String (In_channel.input_all in_)

let init () =
  Runtime.Spec.(define "io"
                (* TODO: mode *)
                +> fun_ "open" Type.Spec.(string @-> stream) "io.open"
                +> fun_ "close" Type.Spec.(string @-> stream) "io.close"
                +> fun_ "read_all" Type.Spec.(stream @-> string) "io.read_all"
                |> end_);
  Runtime.Primitive.add "io.open" prim_open;
  Runtime.Primitive.add "io.read_all" prim_read_all

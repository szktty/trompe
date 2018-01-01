open Base

type t = Runtime.env

let add m ~name ~value =
  Map.add m ~key:name ~data:value

open Base

type t = string

let sep = '.'

let concat root base =
  root ^ (String.of_char sep) ^ base

let base path =
  String.split path ~on:sep
  |> List.last_exn

let parent path =
  String.split path ~on:sep
  |> List.hd

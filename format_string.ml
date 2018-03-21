type placeholder =
  | Text of char
  | Int
  | Float
  | String

type t = placeholder list

let parse s =
  let rec parse0 cs accu =
    match cs with
    | [] -> List.rev accu
    | '%' :: cs -> parse_field cs accu
    | c :: cs -> parse0 cs (Text c :: accu)
  and parse_field cs accu =
    match cs with
    | [] -> List.rev (Text '%' :: accu)
    | 'd' :: cs -> parse0 cs (Int :: accu)
    | 's' :: cs -> parse0 cs (String :: accu)
    | '%' :: cs -> parse0 cs (Text '%' :: accu)
    | c :: cs -> parse0 cs (Text c :: Text '%' :: accu)
  in
  parse0 (String.to_list s) []

let params fmt =
  List.filter fmt ~f:(fun place ->
      match place with
      | Text _ -> false
      | _ -> true)

open Core.Std

let sentencecase word =
  let buf = Buffer.create (String.length word) in
  String.foldi word ~init:()
    ~f:(fun i _ c ->
        match i with
        | 0 -> Buffer.add_char buf c
        | 1 -> Buffer.add_char buf @@ Char.lowercase c
        | _ ->
          if Char.is_uppercase c then begin
            Buffer.add_char buf '_';
          end;
          Buffer.add_char buf @@ Char.lowercase c);
  Buffer.contents buf

module Format = struct

  type placeholder =
    | Text of char
    | Int
    | Float
    | String

  type t = placeholder list

  let parse s =
    let rec parse_format cs accu =
      match cs with
      | [] -> List.rev accu
      | '%' :: cs -> parse_field cs accu
      | c :: cs -> parse_format cs (Text c :: accu)
    and parse_field cs accu =
      match cs with
      | [] -> List.rev (Text '%' :: accu)
      | 'd' :: cs -> parse_format cs (Int :: accu)
      | 's' :: cs -> parse_format cs (String :: accu)
      | '%' :: cs -> parse_format cs (Text '%' :: accu)
      | c :: cs -> parse_format cs (Text c :: Text '%' :: accu)
    in
    parse_format (String.to_list s) []

  let params fmt =
    List.filter fmt ~f:(fun place ->
        match place with
        | Text _ -> false
        | _ -> true)

end

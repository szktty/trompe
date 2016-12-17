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

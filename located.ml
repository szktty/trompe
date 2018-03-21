type 'a t = {
  desc : 'a;
  loc : Location.t option;
}

let create ?loc (desc:'a) : 'a t =
  { desc; loc }

let with_range start_loc end_loc desc =
  create desc ~loc:(Location.union start_loc end_loc)

let with_range_exn start_loc end_loc desc =
  match (start_loc, end_loc) with
  | (Some start_loc, Some end_loc) -> with_range start_loc end_loc desc
  | _ -> failwith "with_range_exn"

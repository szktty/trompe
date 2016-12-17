type 'a t = {
  desc : 'a;
  loc : Location.t option;
}

let create loc desc =
  { desc; loc }

let locate loc desc = create (Some loc) desc

let less desc = create None desc

let with_range start_loc end_loc desc =
  create (Some (Location.union start_loc end_loc)) desc

let with_range_exn start_loc end_loc desc =
  match (start_loc, end_loc) with
  | (Some start_loc, Some end_loc) -> with_range start_loc end_loc desc
  | _ -> failwith "with_range_exn"

type t = {
  name : string;
  sub : t option;
}

let create name =
  { name; sub = None }

let rec add path name =
  match path.sub with
  | None -> { path with sub = Some (create name) }
  | Some sub -> { path with sub = Some (add sub name) }

let rec iter path ~f =
  f path.name;
  match path.sub with
  | None -> ()
  | Some sub -> iter sub ~f


open Core.Std

type 'a t = {
  parent : 'a t option;
  name : string;
  mutable env : 'a Env.t;
  mutable submods : 'a t list;
}

let create ?parent ?(submods=[]) name env =
  { parent = parent;
    name = name;
    env = env;
    submods = [];
  }

let rec find_submodule ?(prefix=[]) m name =
  match prefix with
  | fst :: rest ->
    begin match find_submodule m fst with
      | Result.Error _ as e -> e
      | Result.Ok sub -> find_submodule ~prefix:rest sub name
    end
  | [] ->
    match List.find m.submods ~f:(fun m -> m.name = name) with
    | None -> Result.Error (m, name)
    | Some sub -> Result.Ok sub

let add_attr m ~name ~value =
  m.env <- Env.add_attr m.env name value

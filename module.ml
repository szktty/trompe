open Core.Std

type t = Lang.module_

let top_modules : t String.Map.t ref = ref String.Map.empty

let register m =
  top_modules := String.Map.add !top_modules ~key:Lang.(m.mod_name) ~data:m

let create ?parent name env =
  { Lang.mod_parent = parent;
    mod_name = name;
    mod_vals = String.Map.empty;
    mod_env = env;
    mod_submods = [];
  }

(* -> (t, (t * string)) Result *)
let rec find_submodule (m : t) ?(prefix=[]) name =
  match prefix with
  | fst :: rest ->
    begin match find_submodule m fst with
      | Result.Error _ as e -> e
      | Result.Ok sub -> find_submodule sub ~prefix:rest name
    end
  | [] ->
    match List.find m.mod_submods ~f:(fun m -> m.mod_name = name) with
    | None -> Result.Error (m, name)
    | Some sub -> Result.Ok sub

let add_value m ~key ~data =
  Lang.{ m with mod_vals = String.Map.add m.mod_vals ~key ~data }

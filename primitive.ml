open Core.Std
open Lang

let g_prims : Lang.primitive String.Map.t ref = ref String.Map.empty

let register prims =
  List.iter prims ~f:(fun (name, f) ->
      g_prims := String.Map.add !g_prims name f)

let find name = String.Map.find !g_prims name

let find_exn name =
  match find name with
  | None -> failwith ("primitive not found: " ^ name)
  | Some f -> f

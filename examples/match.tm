let match_bool x =
  match x with
  | true -> "bool: 1"
  | false -> "bool: 2"
  done

let match_int x =
  match x with
  | 0 -> "int: 1"
  | 1 -> "int: 2"
  | _ -> "int: 3"
  done

let match_string x =
  match x with
  | "" -> "string: 1"
  | "hello" -> "string: 2"
  | s -> "string: 3"
  done

let match_tuple x =
  match x with
  | (true, true) -> "tuple 1"
  | (true, false) -> "tuple 2"
  | (false, true) -> "tuple 3"
  | (false, false) -> "tuple 4"
  done

let match_list x =
  match x with
  | [] -> "list: 1"
  | [1; 2; 3] -> "list: 2"
  | hd :: tl -> "list: 3"
  done

;; show $ match_bool true
;; show $ match_bool false
;; show $ match_int 0
;; show $ match_int 1
;; show $ match_int 2
;; show $ match_string ""
;; show $ match_string "hello"
;; show $ match_string "world"
;; show $ match_list []
;; show $ match_list [1; 2; 3]
;; show $ match_list [1; 2; 3; 4; 5]
;; show $ match_tuple (true, true)
;; show $ match_tuple (true, false)
;; show $ match_tuple (false, true)
;; show $ match_tuple (false, false)

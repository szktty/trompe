let rec sum l =
  match l with
  | [] -> 0
  | hd :: tl -> hd + sum tl
  done

;; print_int $ sum [1; 2; 3; 4; 5; 6; 7; 8; 9; 10]

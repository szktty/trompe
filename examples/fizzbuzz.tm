let fizzbuzz i =
    match (i mod 3, i mod 5) with
    | (0, 0) -> "fizzbuzz"
    | (_, 0) -> "buzz"
    | (0, _) -> "fizz"
    | (_, _) -> String.of_int i

;; for i = 1 to 15 do show $ fizzbuzz i done

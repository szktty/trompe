# based on http://rosettacode.org/wiki/Levenshtein_distance#A_recursive_functional_version

let levenshtein s t =
    let rec dist i j =
        match (i,j) with
        | (i,0) -> i
        | (0,j) -> j
        | (i,j) ->
        if String.get s (i-1) = String.get t (j-1) then
            dist (i-1) (j-1)
        else
            let (d1, d2, d3) = (dist (i-1) j, dist i (j-1), dist (i-1) (j-1)) in
            1 + min d1 (min d2 d3)
    in
    dist (String.length s) (String.length t)

let test s t =
    printf "%s -> %s = %d\n" s t (levenshtein s t)

# kitten -> sitting = 3
# rosettacode -> raisethysword = 8
#;; test "kitten" "sitting"

# warning! too slow!
;; test "rosettacode" "raisethysword"

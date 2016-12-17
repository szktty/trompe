def levenshtein(s: String, t: String)
  def dist(i, j)
    case (i, j) do
    | (i, 0) -> i
    | (0, j) -> j
    | (i, j) ->
      if s[i-1] == t[i-1] then
        dist(i-1, j-1)
      else
        let (d1, d2, d3) = (dist(i-1, j), dist(i, j-1), dist(i-1, j-1))
        1 + min(d1, min(d2, d3))
      end
    end
  end
  dist(s.length(), t.length())
end

def test(s, t)
  printf("%s -> %s = %d\n", s, t, levenshtein(s, t))
end

# kitten -> sitting = 3
# rosettacode -> raisethysword = 8
test("kitten", "sitting")
test("rosettacode", "raisethysword")

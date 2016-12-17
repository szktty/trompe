def fizzbuzz(i) 
  case (i % 3, i % 5) do
    | (0, 0) -> "fizzbuzz"
    | (_, 0) -> "buzz"
    | (0, _) -> "fizz"
    | (_, _) -> Int.to_string(i)
  end
end

for i in 1..15 do
  show(fizzbuzz(i))
end

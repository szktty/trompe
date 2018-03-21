def fizzbuzz(i) 
  case (i % 3, i % 5) of
  when (0, 0) then "fizzbuzz"
  when (_, 0) then "buzz"
  when (0, _) then "fizz"
  when (_, _) then Int.to_string(i)
  end
end

for i in 1..15 do
  show(fizzbuzz(i))
end

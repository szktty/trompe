def fib(n)
  if n <= 1 then
    return n
  else 
    return fib(n-1) + fib(n-2)
  end
end

show(fib(10))

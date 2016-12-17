def hanoi(n, a, b, c)
  if n != 0 then
    hanoi((n - 1), a, c, b)
    printf("Move disk from pole %d to pole %d\n", a, b)
    hanoi((n - 1), c, b, a)
  end
end

hanoi(4, 1, 2, 3)

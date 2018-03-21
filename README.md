# Trompe

Trompe is a strongly-typed scripting language with type inference.
This is developed for handy scripting to clean up chores, which enables quick startup, interpretation (no need build configuration), and detecting type mismatch errors easily and quickly.

Current version is pre-pre-pre-alpha. See examples/ for more detail.

## License

Trompe is licensed under the Apache License, Version 2.0.

## Requirements

- Go 1.10+
- Antlr 4.7.1+

## Build

```
$ make
```

## Grammar

### Comments

```
-- comment end of line
```

### Unit

```
()
```

### Boolean

```
true
false
```

### Integers

```
12345
```

### Floating-Point Numbers

```
123.45
0e10
```

### Strings

```
"hello, world!"
```

### Lists

```
[]
[1, 2, 3]
```

### Tuple

```
(1, 2, 3)
```

### Closure

### Calling Functions

```
f()
f(1, 2 , 3)
```

### Block

```
do
  ...
end
```

### Variable Bindings

```
let x = 1
```

### Defining Functions

```
def f(x) 
  x + 1
end
```

### Conditions

```
if n == 0 then
  show("zero")
else 
  show("other")
end
```

### Loop

```
for i in 1..15 do
  show(i)
end
```

### Pattern Matching

```
case i do
when 0 then "0"
when 1 then "1"
when 2 then "2"
when _ then "_"
end
```

### Type Annotations

# TODO

- Library
- Partial application
- Records
- Variants
- References and dereferences
- Operator definition
- Exception handling
- Modules and traits
- Tail call optimization

## Author

SUZUKI Tetsuya (tetsuya.suzuki@gmail.com)

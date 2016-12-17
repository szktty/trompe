# Trompe

Trompe is a strongly-typed scripting language with type inference.
This is developed for handy scripting to clean up chores, which enables quick startup, interpretation (no need build configuration), and detecting type mismatch errors easily and quickly.

Current version is pre-pre-pre-alpha. See examples/ for more detail.

## License

Trompe is licensed under the Apache License, Version 2.0.

## Requirements

- OCaml 4.04+
- OPAM 1.2.2
- Core
- OMake 0.9.8.6-0.rc1
- Menhir

## Installation

Do at directory repository toplevel.

```
$ opam pin add omake 0.9.8.6-0.rc1
$ opam pin add trompe .
```

## Grammar

### Comments

```
# comment end of line
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
  | 0 -> "0"
  | 1 -> "1"
  | 2 -> "2"
  | _ -> "_"
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

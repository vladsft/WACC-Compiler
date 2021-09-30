# Dynamic Arrays

Dynamic arrays are handled by the `make` expression. `make` allocates an array of type t with length l on the heap.
The free statment has also been updated to allow it to free arrays and pairs.

## Syntax

`<type>[] ident = make(<type>, <length>)`
`free (<type>[] | pair)`

## Semantics
The type of a make expression is an array containing the specified type.

## Code Generation

`make(<type>, <length>)` roughly translates to:
```
mov r0, #length
mul r0, r0, r12
add r0, r0, #4
bl .malloc
str #length, [r0]
```

The free statment needs no change at this stage.

# Enhanced assignments

Enhanced assignments are triggered by the following operators: "+=", "-=", "*=", "/=", "%=", "++", "--".
This operation is appliable to any assignlhs-type identifier such as: normal identifier, array element, pair elem, and pair type.
The "+=", "-=", "*=", "/=", "%=" also include an assignrhs-type identifier to the right with the operand, which is part of the effect that the LHS identifier supports.

## Syntax

Generic syntax:
lhs_ident operator rhs_ident

Equivalency examples:
`a += 3` is equivalent to `a = a + 3`
`i++` is equivalent to `i = i + 1`
`a[i] /= b.fst` is equivalent to `a[i] = a[i] / b`

## Semantics
The type of the modified identifier is the same as declared previously. An identifier needs to be declared prior to this operation.
The semantic check ensures that both the assignlhs and the assignrhs operator are of the type integer, being aware that only numbers are compatible with arithmetic operands.

## Code Generation

`a += 3` roughly translates to:
```
str r4, [sp]
ldr r4, [sp]
ldr r5, =54
adds r4, r4, r5
bl p_check_int_overflow
```

The branch to the overflow check is part of the addition assignment.

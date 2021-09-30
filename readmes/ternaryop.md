# Ternary operator

The ternary operator is a form of syntactic sugar for returning an expression depending on the truth of a certain condition.
An expression `a ? b : c` returns b if a is true, and c otherwise.

## Syntax

Generic syntax:
condition ? value_if_true : value_if_false

Example syntax:
`println i % 2 == 0 ? "even" : "odd"` is equivalent to:
 ___________________
| if i % 2 == 0     |
| then              |           
|    println "even" |                            
| else              |
|    println "odd"  |
| fi                |
|___________________|

## Semantics
The semantic checker ensures that the first expression is a boolean condition and the following expressions have the same type. Therefore, it is erronious to write:

`i % 2 == 0 ? true : "give me an even number"`

because of type differences.

## Code Generation

`println i % 2 == 0 ? "even" : "odd"` translates to:

msg_1:
	.word 4
	.ascii "even"
msg_2:
	.word 3
	.ascii "odd"
main:
    ...
	str r4, [sp]
	ldr r4, [sp]
	ldr r5, =2
	mov r0, r4
	mov r1, r5
	bl p_check_divide_by_zero
	bl __aeabi_idivmod
	mov r4, r1
	ldr r5, =0
	cmp r4, r5
	movEQ r4, #1
	movNE r4, #0
	cmp r4, #0
	bEQ else_0

	ldr r4, =msg_1
	bAL end_0
    ...
else_0:

	ldr r4, =msg_2
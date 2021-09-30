# For Loops

For loops integrate a common use case for while loops into a much more compact, readable syntax. The syntactic structure can be found in the "syntax" section below.
The for loop has four main components:
- The new assignment of the iterator
- The looping condition
- The changing condition of the iterator (note that side effecting instructions such "i++" are applicable, since they are implemented in other extensions here)

## Syntax

Generic syntax:
for (initial; condition; change) do
	statements
done

Example syntax:
_________________________________________
| for (int i = 1; i <= 10; i = i + 1) do |
|   int j = i * 2;                       |
|   println j                            |
| done                                   | 
|_______________________________________ |


## Semantics
The semantic check ensures that the new assignment, the condition and the changing statement are all valid. Note that the initial statement can only be a new assignment here, therefore the following program is invalid:
```
int i;
for (i = 1; i <= 10; i = i + 1) do
    println i
done
```
The semantic check also creates a new scope for the statements that are in the for loop. The variables declared within it are separate from the ones outside the loop, and are deleted once the loop ends. This is done through the creation of a new symbol table in the Check() method.

## Code Generation
 ______________________________________
| for (int i = 0; i < 7; i = i + 1) do |
|   skip                               |
| done                                 | 
|______________________________________|

translates to:

```
initial_for_0:
	ldr r4, =0
	str r4, [sp]
cond_for_0:
	ldr r4, [sp]
	ldr r5, =7
	cmp r4, r5
	movLT r4, #1
	movGE r4, #0
	cmp r4, #1
	bNE for_end_for_0

	ldr r4, [sp]
	ldr r5, =1
	adds r4, r4, r5
	bl p_check_int_overflow

	str r4, [sp]
	bAL cond_for_0
for_end_for_0:
	adds sp, sp, #4
	mov r0, #0
	bl exit
	pop {pc}
```
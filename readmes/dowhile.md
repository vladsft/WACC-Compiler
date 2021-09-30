# Do while

The do while statement creates a loop that executes a specified statement until the test condition evaluates to false. The condition is evaluated after executing the statement, resulting in the specified statement executing at least once.

## Syntax

Generic syntax:
do 
	statements
while condition
done

Example syntax:
_______________
| do            |
|   println i ; |
|   i = i + 1   | 
| while i <= 10 |
| done          |
| _____________ |

## Semantics
The semantic check ensures that the statements and the condition are all valid.
The semantic check creates a new scope for the statements that are in the do while loop. The variables declared within it are separate from the ones outside the loop, and are deleted once the loop ends. This is done through the creation of a new symbol table in the Check() method.

## Code generation
 _______________
| do            |
|   println i ; |
|   i = i + 1   | 
| while i <= 10 |
| done          |
| _____________ |

translates to: 

cond_do_while_0:
	ldr r4, [sp]
	mov r0, r4
	bl p_print_int
	bl p_println
	ldr r4, [sp]
	ldr r5, =1
	adds r4, r4, r5
	bl p_check_int_overflow


	str r4, [sp]
	ldr r4, [sp]
	ldr r5, =10
	cmp r4, r5
	movLE r4, #1
	movGT r4, #0

	cmp r4, #1
	bNE do_while_end_do_while_0
	bAL cond_do_while_0
do_while_end_do_while_0:
	adds sp, sp, #4
	mov r0, #0
	bl exit
	pop {pc}

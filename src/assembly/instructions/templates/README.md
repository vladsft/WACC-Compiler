# Delimiters
<> : fields of internal assembly representation 

# Tags 
<op> : used for operands
<1line> </1line> : used for operands
<$reg> : used for registers

# Fields
<.T> : push or pop in stackInstr, depending on the instruction 
<.ID> : identifier which stands for variable name

# Rules
-> : passes the template to the definition
=> : passes the type alias of the registers to the definition of registers

# Function examples

```
func (m mult) Arm11() string {

}
Uses the multiply operand on two values from source registers and stores the result in the destination register 

```
func (d decrementStack) Arm11() string {

}
Decreases the pointer of the stack to add incoming values on top of it.

```
func (b branch) Arm11() string {
	
}
Sends the control flow of the program to a specified label.

```
func (s storeHeap) Arm11() string {

}
Uses the heap for storing a number of bytes onto it. It starts from the bottom of the stack.
If the offset is not 0, it is specified.
If the multiplier is not 0, it is specified.
It's corresponded to %rbp, the heap pointer.

```
 func (l load) Arm11() {

 }   
Loads the memory address of src via label adr_src into dest.

# Template parameters
sp : stack pointer
r0 : register number 0
#1: immediate value 1
#0: immediate value 0
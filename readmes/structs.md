# Structs

Structures are declared at the top of the file. The structures have fields which can be accessed using the dot (`.`) operator. Structs can also have other structs as fields.

## Syntax

### Declaration

```
struct <ident> is 
   int a
   char b 
   bool c
end
```

### Initialisation
```
<structName> <ident> = <structName>{1, 'a', true }
```

### Field Accesses
```
<structName>.<fieldName>
```

## Semantics

### Declaration
Struct Declarations Nodes are stored in the man program node. The fields can **not** be initialised while declaring the constructor.

### Initialisation
Struct Initialisations are part of the literal AST node. Struct objects need to have at least one feild initialised using the constructor operator (`{}`).

### Field Access
Field Accesses are expressions.

## Code Generation

### Declaration
No assembly code is generated during this phase. The symbol table is updated with the correct offsets for the feilds.

### Initialisation
The objects are declared on the heap.

```
mov r0, <structsize>
bl malloc
str r0, [r0]

ldr r4, =<field_value>
strb r4, [r0]

ldr r4, =<field_value>
str r4, [r0, #<field_offset>]
```

### Field Access
```
str r4, [sp, #<object_offset>]
ldr r4, [sp]
ldrb r4, [r4, #<field_access>]
```

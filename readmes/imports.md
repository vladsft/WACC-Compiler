# Imports

Imports occur before the `begin` statement in a file.

Imports are processed during AST generation. Only functions which are actually called are imported into the AST, in order to minimise the size of the generated binary.

When visiting a wacc program, we first visit each of its imports recursively, then we visit its functions and finally statements. On each function definition, we send the function down a channel. Similarly, when we visit a function call we send the function's name down a separate channel. The visitor multiplexes the 2 channels to filter out all functions which are declared but not used. This approach also makes it easier to switch to concurrent imports.

## Syntax

There are 2 ways to import a file in the WACC language:

1. `import <path to library>`
2. `import <path to library> as <alias>`

## Semantics

Imports are not part of the AST and thus have no bearing on the semantic analysis stage.

## Code Generation

Imported functions are indistinguishable from "native" functions at this stage, so the code generation stage doesn't need to be changed.


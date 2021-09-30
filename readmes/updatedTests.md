1. tests/chunk_00/syntaxErr/noBody.wacc

```wacc
begin end
```
2. tests/chunk_12/syntaxErr/noBodyAfterFuncs.wacc 

```
begin
  int f() is
  return 0 
  end
end
```	

We now allow an empty body.


3. /tests/chunk_10/valid/minusMinusExpr.wacc, /tests/chunk_10/valid/plusPlusExpr.wacc

These tests were updated to now have bracketed expressions. Without the brackets it is a syntax error as the parser identifies the decrement operator.

## Before
```
begin
  println 1--2
end
```
## After
```
begin
  println 1-(-2)
end
```

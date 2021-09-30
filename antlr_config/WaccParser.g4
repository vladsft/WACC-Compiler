parser grammar WaccParser;

options {
    tokenVocab = WaccLexer;
}

program: importfile* BEGIN userType* function* stat? END EOF;

importfile:
    IMPORT waccfile SEMICOLON
    | IMPORT waccfile AS ident SEMICOLON
    ;

waccfile: STRING_LITER;

ident: IDENT;
fieldident: ident (DOT fieldident)?;
libident: fieldident | ident ACCESSOR fieldident;

userType:
    STRUCT ident IS declaration+ END
    | CLASS ident IS declaration+ function* END;

function:
    wacctype ident LPAREN paramlist? RPAREN IS funcbody END;

returnable: (RETURN | EXIT) right = expr;

funcbody:
    (stat SEMICOLON)? (
        returnable
        | IF expr THEN funcbody ELSE funcbody ENDIF
    );

paramlist: param (COMMA param)*;

param: 
    wacctype ident #paramNormal
    | libident ident  #paramUserType;

stat:
    SKIPP                                  # statSkip
    | newassign                            # statNewassign
    | declaration                          # statDeclaration 
    | assign                               # statAssign
    | READ assignlhs                       # statRead
    | FREE expr                            # statFree
    | ACQUIRE fieldident                   # statAcquire
    | RELEASE fieldident                   # statRelease
    | UP fieldident                        # statUp
    | DOWN fieldident                      # statDown
    | RETURN expr                          # statReturn
    | EXIT expr                            # statExit
    | PRINT expr                           # statPrint
    | PRINTLN expr                         # statPrintln
    | fieldident (INC | DEC)               # statTrailingUnOp
    | assignlhs (ENH_PLUS | ENH_MINUS | ENH_STAR | ENH_DIV | ENH_MOD) assignrhs
                                           # statEnhancedAssign
    | IF expr THEN stat ELSE stat ENDIF    # statIf
    | WHILE expr DO stat DONE              # statWhile
    | DO stat WHILE expr DONE              # statDoWhile
    | FOR LPAREN newassign SEMICOLON expr SEMICOLON assign RPAREN DO stat DONE 
                                           # statFor
    | BEGIN stat END                       # statBegin
    | stat SEMICOLON stat                  # statMultiple
    | WACC libident LPAREN arglist? RPAREN # statWacc
    ;

declaration: wacctype ident (ASSIGN assignrhs)?;
newassign: wacctype ident ASSIGN assignrhs;

assign: setlhs assignrhs; 
assignlhs:
    fieldident  # leftIdent
    | arrayelem # leftArrayElem
    | pairtype  # leftPairType
    | pairelem  # leftPairElem;

//This was needed so that badComment.wacc would expect both '[' and '='
setlhs: (libident (LBRACKET expr RBRACKET)* | pairtype | pairelem) ASSIGN;

assignrhs:   
    expr                                     # rightExpr
    | NEWPAIR LPAREN expr COMMA expr RPAREN  # rightNewPair
    | pairelem                               # rightPairElem
    | arrayliter                             # rightArrayLiter
    | libident LBRACES arglist? RBRACES      # rightNewUserType
    | CALL libident LPAREN arglist? RPAREN   # rightFunctionCall
    | MAKE LPAREN wacctype COMMA expr RPAREN # make;

arglist: expr (COMMA expr)*;

pairelem: (FST | SND) right = expr;

wacctype: basetype | arraytype | pairtype | libident;

basetype: INT | BOOL | CHAR | STRING | LOCK | SEMA;

arraytype: (pairtype | basetype) (LBRACKET RBRACKET)+;

pairtype: PAIR LPAREN pairelemtype COMMA pairelemtype RPAREN;

pairelemtype: basetype | arraytype | PAIR;

expr:
    intliter {
        x, err := strconv.Atoi($intliter.text)
        antlrPanic := func(s string) {
            panic(antlr.NewFailedPredicateException(p, "", s))
        }
        if err != nil {
            antlrPanic("Somehow a non integer has been parsed as an integer literal")
        }
        if int(int32(x)) != x {
            antlrPanic(fmt.Sprintf("Integer %d is too large", x))
        }
    }                  # exprIntLiter
    | BOOL_LITER       # exprBoolLiter
    | STRING_LITER     # exprStringLiter
    | pairliter        # exprPairLiter
    | fieldident       # exprIdent
    | arrayelem        # exprArrayElem
    | semaliter        # exprSemaLiter
    | CHAR_LITER       # exprCharLiter
    | unaryoper expr   # exprUnaryOp
    | <assoc=left> left = expr op = (STAR | DIV | MOD) right = expr  # exprBinop
    | <assoc=left> left = expr op = (PLUS | MINUS) right = expr      # exprBinop
    | <assoc=left> left = expr op = (
        GREATER
        | GREATER_OR_EQUAL
        | LESS
        | LESS_OR_EQUAL
    ) right = expr                                                   # exprBinop
    | <assoc=left> left = expr op = (EQUAL | NOT_EQUAL) right = expr # exprBinop
    | <assoc=left> left = expr op = AND right = expr                 # exprBinop
    | <assoc=left> left = expr op = OR right = expr                  # exprBinop
    | expr QMARK expr COLON expr                                     # exprTernaryOp
    | LPAREN expr RPAREN                                             # exprBracketed;

unaryoper: NOT | LEN | ORD | CHR | MINUS | TRYLOCK;
trailingUnOper: INC | DEC;

arrayelem: fieldident (LBRACKET expr RBRACKET)+;

intliter: (PLUS | MINUS)? INT_LITER;

semaliter: SEMA LPAREN intliter RPAREN;

arrayliter: LBRACKET (expr (COMMA expr)*)? RBRACKET;

pairliter: NULL;

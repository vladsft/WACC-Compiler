lexer grammar WaccLexer;

//keywords
BEGIN: 'begin';
END: 'end';
IS: 'is';
SKIPP: 'skip';
FREE: 'free';
RETURN: 'return';
EXIT: 'exit';
IF: 'if';
THEN: 'then';
ELSE: 'else';
ENDIF: 'fi';
WHILE: 'while';
DO: 'do';
FOR: 'for';
DONE: 'done';
CALL: 'call';
WACC: 'wacc';

//lock keywords
ACQUIRE: 'acquire';
RELEASE: 'release';
TRYLOCK: 'try_lock';

//semaphore keywords
UP: 'sema_up';
DOWN: 'sema_down';

//base types
INT: 'int';
BOOL: 'bool';
CHAR: 'char';
STRING: 'string';
LOCK: 'lock';
SEMA: 'sema';

//pair types
PAIR: 'pair';
NEWPAIR: 'newpair';
FST: 'fst';
SND: 'snd';

// I/O
READ: 'read';
PRINT: 'print';
PRINTLN: 'println';
MAKE: 'make';

//un_operators
NOT: '!';
LEN: 'len';
ORD: 'ord';
CHR: 'chr';
INC: '++';
DEC: '--';

//bin_operators
PLUS: '+'; // also and int sign
MINUS: '-'; // also a un_operator and an int sign
ASSIGN: '=';
STAR: '*';
DIV: '/';
MOD: '%';
GREATER: '>';
GREATER_OR_EQUAL: '>=';
LESS: '<';
LESS_OR_EQUAL: '<=';
EQUAL: '==';
NOT_EQUAL: '!=';
AND: '&&';
OR: '||';
ENH_PLUS: '+=';
ENH_MINUS: '-=';
ENH_STAR: '*=';
ENH_DIV: '/=';
ENH_MOD: '%=';

//brackets
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';
LBRACES: '{';
RBRACES: '}';

//User Types
CLASS: 'class';
STRUCT: 'struct';
DOT: '.';

QMARK: '?';
COLON: ':';
SEMICOLON: ';';
COMMA: ',';
fragment UNDERSCORE: '_';

//numbers
fragment DIGIT: '0'..'9';
fragment LOWER_LETTERS: [a-z];
fragment UPPER_LETTERS: [A-Z];
fragment LETTERS: LOWER_LETTERS | UPPER_LETTERS;
fragment CHARACTERS: ~('\'' | '"' | '\\') | ESCAPED_CHARS;

//booleans
fragment TRUE: 'true';
fragment FALSE: 'false';
NULL: 'null';

SINGLE_QUOTE: '\'' ;
DOUBLE_QUOTE: '"' ;

//literals
BOOL_LITER: TRUE | FALSE;
INT_LITER: DIGIT+;
CHAR_LITER: SINGLE_QUOTE CHARACTERS SINGLE_QUOTE;
STRING_LITER: DOUBLE_QUOTE CHARACTERS*? DOUBLE_QUOTE;

ESCAPED_CHARS: '\\' [0btnfr"'\\];

//comments
COMMENT: '#' .*? '\n' -> skip;
//whitespace, whitelines
WHITESPACE: [\t\n \r] -> skip;

//imports
AS: 'as';
IMPORT: 'import';
ACCESSOR: '::';
IDENT: (LETTERS | UNDERSCORE) (LETTERS | DIGIT | UNDERSCORE)*;

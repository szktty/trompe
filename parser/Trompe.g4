grammar Trompe;

chunk
    : block EOF
    ;

block
    : stat* retstat?
    ;

stat
    : ';'
    | letdecl
    | fundef
    | funcall
    | doblock
    | if_
    | case_
    ;

retstat
    : 'return' exp? ';'?
    ;

doblock
    : 'do' block 'end'
    ;

letdecl
    : 'let' pattern '=' exp
    ;

fundef
    : 'def' NAME '(' parlist? ')' block 'end'
    | 'def' NAME '(' parlist? ')' '=' exp
    ;

parlist
    : NAME (',' NAME)*
    ;

if_
    : 'if' exp 'then' block ('elseif' exp 'then' block)* ('else' block)? 'end'
    ;

case_
    : 'case' exp 'of' caseclau? ('else' block)? 'end'
    ;

caseclau
    : 'when' pattern guard? 'then' block
    ;

guard
    : 'in' exp
    ;

pattern
    : unit
    | bool_
    | int_
    | float_
    | string_
    | '[' patlist? ']'
    | '(' patlist? ')'
    ;

patlist
    : pattern (',' pattern)*
    ;

explist
    : exp (',' exp)*
    ;

exp
    : unit
    | bool_
    | int_
    | hexint
    | float_
    | hexfloat
    | string_
    | prefixexp
    | list
    | tuple
    | tableconstructor
    | anonfun
    | statexp
    | <assoc=right> exp operatorPower exp
    | operatorUnary exp
    | exp operatorMulDivMod exp
    | exp operatorAddSub exp
    | <assoc=right> exp operatorStrcat exp
    | exp operatorComparison exp
    | exp operatorAnd exp
    | exp operatorOr exp
    | exp operatorBitwise exp
    ;

prefixexp
    : var_OrExp nameAndArgs*
    ;

funcall
    : var_OrExp arglist
    ;

var_OrExp
    : var_ | '(' exp ')'
    ;

var_
    : (modulepath | '(' exp ')' var_Suffix) var_Suffix*
    ;

modulepath
    : NAME ('.' NAME)*
    ;

var_Suffix
    : nameAndArgs* ('[' exp ']' | '.' NAME)
    ;

nameAndArgs
    : (':' NAME)? arglist
    ;

/*
var_
    : NAME | prefixexp '[' exp ']' | prefixexp '.' NAME
    ;
prefixexp
    : var_ | funcall | '(' exp ')'
    ;
funcall
    : prefixexp args | prefixexp ':' NAME args
    ;
*/

arglist
    : '(' explist? ')'
    ;

list
    : '[' eltlist? ']'
    ;

tuple
    : '(' eltlist? ')'
    ;

eltlist
    : exp (',' exp)* ','?
    ;

tableconstructor
    : '{' fieldlist? '}'
    ;

fieldlist
    : field (',' field)* ','?
    ;

field
    : NAME '=' exp
    ;

anonfun
    : '[' (parlist | unit)? 'in' block ']'
    ;

statexp
    : '[' stat ']'
    ;

operatorOr
	: 'or';

operatorAnd
	: 'and';

operatorComparison
	: '<' | '>' | '<=' | '>=' | '~=' | '==';

operatorStrcat
	: '..';

operatorAddSub
	: '+' | '-';

operatorMulDivMod
	: '*' | '/' | '%' | '//';

operatorBitwise
	: '&' | '|' | '~' | '<<' | '>>';

operatorUnary
    : '#' | '-' | '~';

operatorPower
    : '^';

unit
    : '(' ')'
    ;

bool_
    : 'true'
    | 'false'
    ;

int_
    : INT
    ;

hexint
    : HEX
    ;

float_
    : FLOAT
    ;

hexfloat
    : HEX_FLOAT
    ;

string_
    : NORMALSTRING | CHARSTRING | LONGSTRING
    ;

// LEXER

NAME
    : [a-zA-Z_][a-zA-Z_0-9]*
    ;

NORMALSTRING
    : '"' ( EscapeSequence | ~('\\'|'"') )* '"'
    ;

CHARSTRING
    : '\'' ( EscapeSequence | ~('\''|'\\') )* '\''
    ;

LONGSTRING
    : '[' NESTED_STR ']'
    ;

fragment
NESTED_STR
    : '=' NESTED_STR '='
    | '[' .*? ']'
    ;

INT
    : Digit+
    ;

HEX
    : '0' [xX] HexDigit+
    ;

FLOAT
    : Digit+ '.' Digit* ExponentPart?
    | '.' Digit+ ExponentPart?
    | Digit+ ExponentPart
    ;

HEX_FLOAT
    : '0' [xX] HexDigit+ '.' HexDigit* HexExponentPart?
    | '0' [xX] '.' HexDigit+ HexExponentPart?
    | '0' [xX] HexDigit+ HexExponentPart
    ;

fragment
ExponentPart
    : [eE] [+-]? Digit+
    ;

fragment
HexExponentPart
    : [pP] [+-]? Digit+
    ;

fragment
EscapeSequence
    : '\\' [abfnrtvz"'\\]
    | '\\' '\r'? '\n'
    | DecimalEscape
    | HexEscape
    | UtfEscape
    ;

fragment
DecimalEscape
    : '\\' Digit
    | '\\' Digit Digit
    | '\\' [0-2] Digit Digit
    ;

fragment
HexEscape
    : '\\' 'x' HexDigit HexDigit
    ;
fragment
UtfEscape
    : '\\' 'u{' HexDigit+ '}'
    ;
fragment
Digit
    : [0-9]
    ;
fragment
HexDigit
    : [0-9a-fA-F]
    ;
COMMENT
    : '--[' NESTED_STR ']' -> channel(HIDDEN)
    ;

LINE_COMMENT
    : '--'
    (                                               // --
    | '[' '='*                                      // --[==
    | '[' '='* ~('='|'['|'\r'|'\n') ~('\r'|'\n')*   // --[==AA
    | ~('['|'\r'|'\n') ~('\r'|'\n')*                // --AAA
    ) ('\r\n'|'\r'|'\n'|EOF)
    -> channel(HIDDEN)
    ;

WS
    : [ \t\u000C\r\n]+ -> skip
    ;

SHEBANG
    : '#' '!' ~('\n'|'\r')* -> channel(HIDDEN)
    ;

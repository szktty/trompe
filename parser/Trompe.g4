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
    | for_
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

for_
    : 'for' pattern 'in' exp 'do' block 'end'
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
    | pattern rangeop pattern
    | '[' patlist? ']'
    | '(' patlist? ')'
    | NAME
    ;

patlist
    : pattern (',' pattern)*
    ;

exp
    : simpleexp
    | funcall
    /*
    | <assoc=right> exp operatorPower exp
    | operatorUnary exp
    | exp operatorMulDivMod exp
    | exp operatorAddSub exp
    | exp operatorComparison exp
    | exp operatorAnd exp
    | exp operatorOr exp
    | exp operatorBitwise exp
    */
    | left=exp rangeop right=exp
    ;

parenexp
    : o='(' exp c=')'
    ;

simpleexp
    : unit
    | bool_
    | int_
    | hexint
    | float_
    | hexfloat
    | string_
    | list
    | tuple
    | tableconstructor
    | anonfun
    | statexp
    | var_
    | parenexp
    ;

funcall
    : simpleexp arglist
    ;

arglist
    : '(' explist? ')'
    ;

explist
    : exp (',' exp)*
    ;

var_
    : modulepath
    ;

modulepath
    : NAME ('.' NAME)*
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
    : '[' (parlist | unit)? 'in' stat* exp ']'
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

rangeop
    : '...'
    | '..<'
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
    : Digit+ Fraction? ExponentPart?
    ;

HEX_FLOAT
    : '0' [xX] HexDigit+ '.' HexDigit* HexExponentPart?
    | '0' [xX] '.' HexDigit+ HexExponentPart?
    | '0' [xX] HexDigit+ HexExponentPart
    ;

fragment
Fraction
    : '.' Digit+
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

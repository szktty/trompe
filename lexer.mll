{
open Core.Std
open Lexing
open Parser

exception Syntax_error of Position.t * string

let next_line lexbuf =
  let pos = lexbuf.lex_curr_p in
  lexbuf.lex_curr_p <-
    { pos with pos_bol = 1;
               pos_lnum = pos.pos_lnum + 1
    }

let revise_pos pos lexbuf =
  Position.of_lexing_pos
    { pos with pos_bol = pos.pos_cnum - lexbuf.lex_curr_p.pos_bol + 1 }

let start_pos lexbuf =
  revise_pos (lexeme_start_p lexbuf) lexbuf

let end_pos lexbuf =
  let p = lexeme_end_p lexbuf in
  let p' = { p with pos_cnum = p.pos_cnum - 1 } in
  revise_pos p' lexbuf

let to_loc lexbuf =
  Location.create (start_pos lexbuf) (end_pos lexbuf)

let to_word lexbuf =
  Located.locate (to_loc lexbuf) (lexeme lexbuf)

let to_word_map lexbuf ~f =
  let value = f @@ lexeme lexbuf in
  Located.locate (to_loc lexbuf) (f @@ lexeme lexbuf)

let strlit lexbuf read =
  let sp = start_pos lexbuf in
  let contents = read (Buffer.create 17) lexbuf in
  let loc = Location.create sp (end_pos lexbuf) in
  Located.locate loc contents

}

let int = ['0'-'9'] ['0'-'9']*
let digit = ['0'-'9']
let frac = '.' digit+
let exp = ['e' 'E'] ['-' '+']? digit+
let float = digit+ frac? exp?
let white = [' ' '\t']+
let newline = '\r' | '\n' | "\r\n"
let lident = ['a'-'z' '_'] ['a'-'z' 'A'-'Z' '0'-'9' '_' '\'']* ['!' '?']?
let uident = ['A'-'Z'] ['a'-'z' 'A'-'Z' '0'-'9' '_' '\'']*
let comment = '#'
let blank = [' ' '\t']*

rule read =
  parse
  | white       { read lexbuf }
  | newline     { next_line lexbuf; read lexbuf }
  | comment     { skip_comment lexbuf; read lexbuf }
  | int         { INT (to_word_map lexbuf ~f:int_of_string) }
  | '-' int     { NEG_INT (to_word_map lexbuf ~f:int_of_string) } (* TODO: neg *)
  | float       { FLOAT (to_word_map lexbuf ~f:float_of_string) }
  | '-' float       { FLOAT (to_word_map lexbuf ~f:float_of_string) } (* TODO: neg *)
  | '"'         { STRING (strlit lexbuf read_string) } 
  | '('         { LPAREN (to_loc lexbuf) }
  | ')'         { RPAREN }
  | '{'         { LBRACE }
  | '}'         { RBRACE }
  | '['         { LBRACK }
  | ']'         { RBRACK }
  | ':'         { COLON }
  | ','         { COMMA }
  | '.'         { DOT }
  | ".."        { DOT2 }
  | '|'         { BAR }
  | '^'         { CARET }
  | "<-"        { LARROW }
  | "->"        { RARROW }
  | '+'         { PLUS (to_loc lexbuf) }
  | '-'         { MINUS (to_loc lexbuf) }
  | '*'         { AST (to_loc lexbuf) }
  | '/'         { SLASH (to_loc lexbuf) }
  | '%'         { PCT (to_loc lexbuf) }
  | "**"        { AST2 (to_loc lexbuf) }
  | '@'         { AT (to_loc lexbuf) }
  | '&'         { AMP (to_loc lexbuf) }
  | '='         { EQ (to_loc lexbuf) }
  | "=="        { EQQ (to_loc lexbuf) }
  | "!="        { NE (to_loc lexbuf) }
  | "<"         { LT (to_loc lexbuf) }
  | "<="        { LE (to_loc lexbuf) }
  | ">"         { GT (to_loc lexbuf) }
  | ">="        { GE (to_loc lexbuf) }
  | "<<"        { LCOMP (to_loc lexbuf) }
  | ">>"        { RCOMP (to_loc lexbuf) }
  | "<|"        { LPIPE (to_loc lexbuf) }
  | "|>"        { RPIPE (to_loc lexbuf) }
  | "and"       { AND (to_loc lexbuf) }
  | "or"        { OR (to_loc lexbuf) }
  | "case"      { CASE }
  | "catch"     { CATCH }
  | "def"       { DEF }
  | "do"        { DO }
  | "else"      { ELSE }
  | "end"       { END }
  | "for"       { FOR }
  | "fun"       { FUN }
  | "if"        { IF }
  | "in"        { IN }
  | "let"       { LET }
  | "module"    { MODULE }
  | "raise"     { RAISE }
  | "return"    { RETURN }
  | "then"      { THEN }
  | "try"       { TRY }
  | "when"      { WHEN }
  | "false"     { FALSE (to_loc lexbuf) }
  | "true"      { TRUE (to_loc lexbuf) }
  | lident      { LIDENT (to_word lexbuf) }
  | '*' lident  { DEREF_LIDENT (to_word_map lexbuf
                     ~f:(fun s -> String.chop_prefix_exn s ~prefix:"*")) }
  | uident
  {
    let ident = lexeme lexbuf in
    let expected = Utils.sentencecase ident in
    if expected <> ident then begin
      raise (Syntax_error (start_pos lexbuf,
         Printf.sprintf "Identifier '%s' must be sentence case with underscore (example: '%s')"
            ident expected))
    end;
    UIDENT (to_word lexbuf)
  }
  | _           { raise (Syntax_error (start_pos lexbuf, "Unexpected char: " ^ Lexing.lexeme lexbuf)) }
  | eof         { EOF }

and skip_comment =
  parse
  | newline     { next_line lexbuf }
  | eof         { () }
  | _           { skip_comment lexbuf }

and read_string buf =
  parse
  | '"'       { Buffer.contents buf }
  | '\\' '/'  { Buffer.add_char buf '/'; read_string buf lexbuf }
  | '\\' '\\' { Buffer.add_char buf '\\'; read_string buf lexbuf }
  | '\\' 'b'  { Buffer.add_char buf '\b'; read_string buf lexbuf }
  | '\\' 'f'  { Buffer.add_char buf '\012'; read_string buf lexbuf }
  | '\\' 'n'  { Buffer.add_char buf '\n'; read_string buf lexbuf }
  | '\\' 'r'  { Buffer.add_char buf '\r'; read_string buf lexbuf }
  | '\\' 't'  { Buffer.add_char buf '\t'; read_string buf lexbuf }
  | [^ '"' '\\']+
    { Buffer.add_string buf (Lexing.lexeme lexbuf);
      read_string buf lexbuf
    }
  | _ { raise (Syntax_error (start_pos lexbuf, "Illegal string character: " ^ Lexing.lexeme lexbuf)) }
  | eof { raise (Syntax_error (start_pos lexbuf, "String is not terminated")) }

" Vim syntax file for Trompe
" This file is based on omlet.vim.
"
" Usage:
" - Copy this file under $HOME/.vim/syntax directory
" - Write the following command into $HOME/.vim/ftdetect/trompe.vim 
"   au BufRead,BufNewFile *.tm,*.tmi set ft=trompe

" ============================================================================
" Original description (omlet.vim)
" ============================================================================
" Vim syntax file
" Language:     OCaml
" Filenames:    *.ml *.mli *.mll *.mly
" Maintainers:  Markus Mottl      <markus.mottl@gmail.com>
"               Karl-Heinz Sylla  <Karl-Heinz.Sylla@gmd.de>
"               Issac Trotts      <ijtrotts@ucdavis.edu>
" URL:          http://www.ocaml.info/vim/syntax/ocaml.vim
" Last Change:  2004 Aug 13 - Added new type keywords (MM)
"               2004 Jul 30 - Added script keyword "thread" (MM)
"               2004 May 15 - Added keyword "format4" (MM)
"               2003 Jan 19 - Added script keyword "require" (MM)

" A minor patch was applied to the official version so that object/end
" can be distinguished from begin/end, which is used for indentation,
" and folding. (David Baelde)
" ============================================================================

" For version 5.x: Clear all syntax items
" For version 6.x: Quit when a syntax file was already loaded
if version < 600
  syntax clear
elseif exists("b:current_syntax") && b:current_syntax != "trompe"
  finish
endif

" Trompe is case sensitive.
syn case match

" Script headers highlighted like comments
syn match    trompeComment   "^#!.*"
syn match    trompeComment   "#.*"

" Scripting directives
syn match    trompeScript "^#\<\(quit\|labels\|warnings\|directory\|cd\|load\|use\|install_printer\|remove_printer\|require\|thread\|trace\|untrace\|untrace_all\|print_depth\|print_length\)\>"

" lowercase identifier - the standard way to match
syn match    trompeLCIdentifier /\<\(\l\|_\)\(\w\|'\)*\>/

syn match    trompeKeyChar    "|"

" Errors
syn match    trompeBraceErr   "}"
syn match    trompeBrackErr   "\]"
syn match    trompeParenErr   ")"
syn match    trompeArrErr     "|]"

syn match    trompeCommentErr "\*)"

syn match    trompeCountErr   "\<downto\>"
syn match    trompeCountErr   "\<to\>"

syn match    trompeDoErr      "\<do\>"

syn match    trompeDoneErr    "\<done\>"
syn match    trompeThenErr    "\<then\>"

" Error-highlighting of "end" without synchronization:
" as keyword or as error (default)
if exists("trompe_noend_error")
  syn match    trompeKeyword    "\<end\>"
else
  syn match    trompeEndErr     "\<end\>"
endif

" Some convenient clusters
syn cluster  trompeAllErrs contains=trompeBraceErr,trompeBrackErr,trompeParenErr,trompeCommentErr,trompeCountErr,trompeDoErr,trompeDoneErr,trompeEndErr,trompeThenErr

syn cluster  trompeAENoParen contains=trompeBraceErr,trompeBrackErr,trompeCommentErr,trompeCountErr,trompeDoErr,trompeDoneErr,trompeEndErr,trompeThenErr

syn cluster  trompeContained contains=trompeTodo,trompePreDef,trompeModParam,trompeModParam1,trompePreMPRestr,trompeMPRestr,trompeMPRestr1,trompeMPRestr2,trompeMPRestr3,trompeModRHS,trompeFuncWith,trompeFuncStruct,trompeModTypeRestr,trompeModTRWith,trompeWith,trompeWithRest,trompeModType,trompeFullMod


" Enclosing delimiters
syn region   trompeEncl transparent matchgroup=trompeKeyword start="(" matchgroup=trompeKeyword end=")" contains=ALLBUT,@trompeContained,trompeParenErr
syn region   trompeEncl transparent matchgroup=trompeKeyword start="{" matchgroup=trompeKeyword end="}"  contains=ALLBUT,@trompeContained,trompeBraceErr
syn region   trompeEncl transparent matchgroup=trompeKeyword start="\[" matchgroup=trompeKeyword end="\]" contains=ALLBUT,@trompeContained,trompeBrackErr
syn region   trompeEncl transparent matchgroup=trompeKeyword start="\[|" matchgroup=trompeKeyword end="|\]" contains=ALLBUT,@trompeContained,trompeArrErr


" Comments
syn region   trompeComment start="(\*" end="\*)" contains=trompeComment,trompeTodo
syn keyword  trompeTodo contained TODO FIXME XXX

" Blocks
syn region   trompeEnd matchgroup=trompeKeyword start="\<begin\>" matchgroup=trompeKeyword end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr

" "for"
syn region   trompeEnd matchgroup=trompeKeyword start="\<for\>" matchgroup=trompeKeyword end="\<done\>" contains=ALLBUT,@trompeContained,trompeCountErr

" "do"
syn region   trompeDo matchgroup=trompeKeyword start="\<do\>" matchgroup=trompeKeyword end="\<done\>" contains=ALLBUT,@trompeContained,trompeDoneErr

" "if"
syn region   trompeEnd matchgroup=trompeKeyword start="\<if\>" matchgroup=trompeKeyword end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr

" "match"
syn region   trompeDo matchgroup=trompeKeyword start="\<match\>" matchgroup=trompeKeyword end="\<done\>" contains=ALLBUT,@trompeContained,trompeDoneErr

" "function"
syn region   trompeEnd matchgroup=trompeKeyword start="\<function\>" matchgroup=trompeKeyword end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr

" "fun"
syn region   trompeEnd matchgroup=trompeKeyword start="\<fun\>" matchgroup=trompeKeyword end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr



"" Modules

" "struct"
syn region   trompeStruct matchgroup=trompeModule start="\<struct\>" matchgroup=trompeModule end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr

" "sig"
syn region   trompeSig matchgroup=trompeModule start="\<sig\>" matchgroup=trompeModule end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr,trompeModule
syn region   trompeModSpec matchgroup=trompeKeyword start="\<module\>" matchgroup=trompeModule end="\<\u\(\w\|'\)*\>" contained contains=@trompeAllErrs,trompeComment skipwhite skipempty nextgroup=trompeModTRWith,trompeMPRestr

" "open"
syn region   trompeStruct matchgroup=trompeModule start="\<open\>" matchgroup=trompeModule end="\<import\>" contains=ALLBUT,@trompeContained,trompeEndErr

" "import"
syn region   trompeNone matchgroup=trompeKeyword start="\<import\>" matchgroup=trompeModule end="\<\u\(\w\|'\)*\(\.\u\(\w\|'\)*\)*\>" contains=@trompeAllErrs,trompeComment

" "include"
syn match    trompeKeyword "\<include\>" contained skipwhite skipempty nextgroup=trompeModParam,trompeFullMod

" "module" - somewhat complicated stuff ;-)
syn region   trompeModule matchgroup=trompeKeyword start="\<module\>" matchgroup=trompeModule end="\<\u\(\w\|'\)*\>" contains=@trompeAllErrs,trompeComment skipwhite skipempty nextgroup=trompePreDef
syn region   trompePreDef start="."me=e-1 matchgroup=trompeKeyword end="\l\|="me=e-1 contained contains=@trompeAllErrs,trompeComment,trompeModParam,trompeModTypeRestr,trompeModTRWith nextgroup=trompeModPreRHS
syn region   trompeModParam start="([^*]" end=")" contained contains=@trompeAENoParen,trompeModParam1
syn match    trompeModParam1 "\<\u\(\w\|'\)*\>" contained skipwhite skipempty nextgroup=trompePreMPRestr

syn region   trompePreMPRestr start="."me=e-1 end=")"me=e-1 contained contains=@trompeAllErrs,trompeComment,trompeMPRestr,trompeModTypeRestr

syn region   trompeMPRestr start=":" end="."me=e-1 contained contains=@trompeComment skipwhite skipempty nextgroup=trompeMPRestr1,trompeMPRestr2,trompeMPRestr3
syn region   trompeMPRestr1 matchgroup=trompeModule start="\ssig\s\=" matchgroup=trompeModule end="\<end\>" contained contains=ALLBUT,@trompeContained,trompeEndErr,trompeModule
syn region   trompeMPRestr2 start="\sfunctor\(\s\|(\)\="me=e-1 matchgroup=trompeKeyword end="->" contained contains=@trompeAllErrs,trompeComment,trompeModParam skipwhite skipempty nextgroup=trompeFuncWith,trompeMPRestr2
syn match    trompeMPRestr3 "\w\(\w\|'\)*\(\.\w\(\w\|'\)*\)*" contained
syn match    trompeModPreRHS "=" contained skipwhite skipempty nextgroup=trompeModParam,trompeFullMod
syn region   trompeModRHS start="." end=".\w\|([^*]"me=e-2 contained contains=trompeComment skipwhite skipempty nextgroup=trompeModParam,trompeFullMod
syn match    trompeFullMod "\<\u\(\w\|'\)*\(\.\u\(\w\|'\)*\)*" contained skipwhite skipempty nextgroup=trompeFuncWith

syn region   trompeFuncWith start="([^*]"me=e-1 end=")" contained contains=trompeComment,trompeWith,trompeFuncStruct skipwhite skipempty nextgroup=trompeFuncWith
syn region   trompeFuncStruct matchgroup=trompeModule start="[^a-zA-Z]struct\>"hs=s+1 matchgroup=trompeModule end="\<end\>" contains=ALLBUT,@trompeContained,trompeEndErr

syn match    trompeModTypeRestr "\<\w\(\w\|'\)*\(\.\w\(\w\|'\)*\)*\>" contained
syn region   trompeModTRWith start=":\s*("hs=s+1 end=")" contained contains=@trompeAENoParen,trompeWith
syn match    trompeWith "\<\(\u\(\w\|'\)*\.\)*\w\(\w\|'\)*\>" contained skipwhite skipempty nextgroup=trompeWithRest
syn region   trompeWithRest start="[^)]" end=")"me=e-1 contained contains=ALLBUT,@trompeContained

" "module type"
syn region   trompeKeyword start="\<module\>\s*\<type\>" matchgroup=trompeModule end="\<\w\(\w\|'\)*\>" contains=trompeComment skipwhite skipempty nextgroup=trompeMTDef
syn match    trompeMTDef "=\s*\w\(\w\|'\)*\>"hs=s+1,me=s

syn keyword  trompeKeyword  and as assert class
syn keyword  trompeKeyword  constraint else then
syn keyword  trompeKeyword  exception external fun function

syn keyword  trompeKeyword  in inherit initializer
syn keyword  trompeKeyword  land lazy let 
syn keyword  trompeKeyword  method mutable new of
syn keyword  trompeKeyword  parser private raise rec
syn keyword  trompeKeyword  try type
syn keyword  trompeKeyword  val virtual when while with

syn keyword  trompeKeyword  do function
syn keyword  trompeBoolean  true false
syn match    trompeKeyChar  "!"

syn keyword  trompeType     array bool char exn float format format4
syn keyword  trompeType     int int32 int64 lazy_t list nativeint option
syn keyword  trompeType     string unit

syn keyword  trompeOperator asr lor lsl lsr lxor mod not

syn match    trompeConstructor  "(\s*)"
syn match    trompeConstructor  "\[\s*\]"
syn match    trompeConstructor  "\[|\s*>|]"
syn match    trompeConstructor  "\[<\s*>\]"
syn match    trompeConstructor  "\u\(\w\|'\)*\>"

" Polymorphic variants
syn match    trompeConstructor  "`\w\(\w\|'\)*\>"

" Module prefix
syn match    trompeModPath      "\u\(\w\|'\)*\."he=e-1

syn match    trompeCharacter    "'\\\d\d\d'\|'\\[\'ntbr]'\|'.'"
syn match    trompeCharErr      "'\\\d\d'\|'\\\d'"
syn match    trompeCharErr      "'\\[^\'ntbr]'"
syn region   trompeString       start=+"+ skip=+\\\\\|\\"+ end=+"+

syn match    trompeFunDef       "->"
syn match    trompeRefAssign    ":="
syn match    trompeTopStop      ";;"
syn match    trompeOperator     "\^"
syn match    trompeOperator     "::"

syn match    trompeOperator     "&&"
syn match    trompeOperator     "<"
syn match    trompeOperator     ">"
syn match    trompeAnyVar       "\<_\>"
syn match    trompeKeyChar      "|[^\]]"me=e-1
syn match    trompeKeyChar      ";"
syn match    trompeKeyChar      "\~"
syn match    trompeKeyChar      "?"
syn match    trompeKeyChar      "\*"
syn match    trompeKeyChar      "="
syn match    trompeOperator     "\$"

syn match    trompeOperator   "<-"

syn match    trompeNumber        "\<-\=\d\+\>"
syn match    trompeNumber        "\<-\=0[x|X]\x\+\>"
syn match    trompeNumber        "\<-\=0[o|O]\o\+\>"
syn match    trompeNumber        "\<-\=0[b|B][01]\+\>"
syn match    trompeFloat         "\<-\=\d\+\.\d*\([eE][-+]\=\d\+\)\=[fl]\=\>"

" Labels
syn match    trompeLabel        "\~\(\l\|_\)\(\w\|'\)*"lc=1
syn match    trompeLabel        "?\(\l\|_\)\(\w\|'\)*"lc=1
syn region   trompeLabel transparent matchgroup=trompeLabel start="?(\(\l\|_\)\(\w\|'\)*"lc=2 end=")"me=e-1 contains=ALLBUT,@trompeContained,trompeParenErr


" Synchronization
syn sync minlines=50
syn sync maxlines=500

syn sync match trompeDoSync      grouphere  trompeDo      "\<do\>"
syn sync match trompeDoSync      groupthere trompeDo      "\<done\>"

if exists("trompe_revised")
    syn sync match trompeEndSync     grouphere  trompeEnd     "\<\(object\)\>"
else
    syn sync match trompeEndSync     grouphere  trompeEnd     "\<\(begin\|object\)\>"
endif

syn sync match trompeEndSync     groupthere trompeEnd     "\<end\>"
syn sync match trompeStructSync  grouphere  trompeStruct  "\<struct\>"
syn sync match trompeStructSync  groupthere trompeStruct  "\<end\>"
syn sync match trompeSigSync     grouphere  trompeSig     "\<sig\>"
syn sync match trompeSigSync     groupthere trompeSig     "\<end\>"

" Define the default highlighting.
" For version 5.7 and earlier: only when not done already
" For version 5.8 and later: only when an item doesn't have highlighting yet
if version >= 508 || !exists("did_trompe_syntax_inits")
  if version < 508
    let did_trompe_syntax_inits = 1
    command -nargs=+ HiLink hi link <args>
  else
    command -nargs=+ HiLink hi def link <args>
  endif

  HiLink trompeBraceErr     Error
  HiLink trompeBrackErr     Error
  HiLink trompeParenErr     Error
  HiLink trompeArrErr       Error

  HiLink trompeCommentErr   Error

  HiLink trompeCountErr     Error
  HiLink trompeDoErr        Error
  HiLink trompeDoneErr      Error
  HiLink trompeEndErr       Error
  HiLink trompeThenErr      Error

  HiLink trompeCharErr      Error

  HiLink trompeErr          Error

  HiLink trompeComment      Comment

  HiLink trompeModPath      Include
  HiLink trompeObject	   Include
  HiLink trompeModule       Include
  HiLink trompeModParam1    Include
  HiLink trompeModType      Include
  HiLink trompeMPRestr3     Include
  HiLink trompeFullMod      Include
  HiLink trompeModTypeRestr Include
  HiLink trompeWith         Include
  HiLink trompeMTDef        Include

  HiLink trompeScript       Include

  HiLink trompeConstructor  Constant

  HiLink trompeModPreRHS    Keyword
  HiLink trompeMPRestr2     Keyword
  HiLink trompeKeyword      Keyword
  HiLink trompeFunDef       Keyword
  HiLink trompeRefAssign    Keyword
  HiLink trompeKeyChar      Keyword
  HiLink trompeAnyVar       Keyword
  HiLink trompeTopStop      Keyword
  HiLink trompeOperator     Keyword

  HiLink trompeBoolean      Boolean
  HiLink trompeCharacter    Character
  HiLink trompeNumber       Number
  HiLink trompeFloat        Float
  HiLink trompeString       String

  HiLink trompeLabel        Identifier

  HiLink trompeType         Type

  HiLink trompeTodo         Todo

  HiLink trompeEncl         Keyword

  delcommand HiLink
endif

let b:current_syntax = "trompe"

" vim: ts=8

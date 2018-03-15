.PHONY: all

all:
	go build -o trompe cmd/trompe/main.go

syntax:
	antlr4 -Dlanguage=Go trompe.g4

clean:
	rm -f trompe.interp trompe.tokens trompeLexer.interp trompeLexer.tokens trompe_base_listener.go trompe_lexer.go trompe_listener.go trompe_parser.go trompe


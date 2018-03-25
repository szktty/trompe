.PHONY: all trompe trompec

all: syntax trompe trompec

trompe:
	go build -o trompe cmd/trompe/main.go

trompec:
	go build -o trompec cmd/trompec/main.go

syntax:
	antlr4 -Dlanguage=Go parser/Trompe.g4

clean:
	rm -f parser/trompe.interp parser/trompe.tokens parser/trompeLexer.interp parser/trompeLexer.tokens parser/trompe_base_listener.go parser/trompe_lexer.go parser/trompe_listener.go parser/trompe_parser.go trompe


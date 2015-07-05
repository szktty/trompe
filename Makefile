FLAGS = -d -v
PROG = trompe

all:
	go tool yacc -o shared/parser.go shared/parser.go.y
	go build

clean:
	rm -f $(PROG) shared/parser.go y.output

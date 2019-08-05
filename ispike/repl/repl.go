package repl

import (
	"bufio"
	"fmt"
	"io"
	"spike-interpreter-go/spike/lexer"
	"strings"
)

const prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		l := lexer.NewLexer(strings.NewReader(scanner.Text()))

		for token, err := l.NextToken(); token.Type != lexer.Eof; token, err = l.NextToken() {
			if err != nil {
				fmt.Print(err)
				return
			}

			fmt.Fprintf(out, "%+v\n", token)
		}
	}
}

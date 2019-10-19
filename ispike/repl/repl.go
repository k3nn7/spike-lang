package repl

import (
	"bufio"
	"fmt"
	"io"
	"spike-interpreter-go/spike/eval"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
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

		l := lexer.New(strings.NewReader(scanner.Text()))
		p := parser.New(l)
		program, err := p.ParseProgram()

		if err != nil {
			fmt.Print(err)
			return
		}

		result, err := eval.Eval(program)
		if err != nil {
			fmt.Print(err)
		}

		fmt.Fprint(out, result.Inspect())
		fmt.Fprint(out, "\n")
	}
}

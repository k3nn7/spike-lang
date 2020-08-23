package repl

import (
	"bufio"
	"fmt"
	"io"
	"spike-interpreter-go/spike/compiler"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
	"spike-interpreter-go/spike/vm"
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

		c := compiler.New()
		err = c.Compile(program)
		if err != nil {
			fmt.Print(err)
			return
		}

		v := vm.New(c.Bytecode())
		err = v.Run()
		if err != nil {
			fmt.Print(err)
			return
		}

		fmt.Fprint(out, v.LastPoppedStackElement().Inspect())
		fmt.Fprint(out, "\n")
	}
}

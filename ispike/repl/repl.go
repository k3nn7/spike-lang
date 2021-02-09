package repl

import (
	"bufio"
	"fmt"
	"io"
	"spike-interpreter-go/spike/compiler"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/object"
	"spike-interpreter-go/spike/parser"
	"spike-interpreter-go/spike/vm"
	"strings"
)

const prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		_, err := fmt.Fprint(out, prompt)
		if err != nil {
			fmt.Print(err)
			return
		}

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

		c := compiler.NewWithState(symbolTable, constants)
		err = c.Compile(program)
		if err != nil {
			fmt.Print(err)
			return
		}

		v := vm.NewWithGlobalStore(c.Bytecode(), globals)
		err = v.Run()
		if err != nil {
			fmt.Print(err)
			return
		}

		_, err = fmt.Fprint(out, v.LastPoppedStackElement().Inspect())
		if err != nil {
			fmt.Print(err)
			return
		}

		_, err = fmt.Fprint(out, "\n")
		if err != nil {
			fmt.Print(err)
			return
		}
	}
}

package main

import (
	"fmt"
	"os"
	"spike-interpreter-go/spike/eval"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
)

func main() {
	l := lexer.New(os.Stdin)
	p := parser.New(l)

	program, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("Parser error: %s\n", err)
		return
	}

	result, err := eval.Eval(program)
	if err != nil {
		fmt.Printf("Runtime error: %s\n", err)
		return
	}

	fmt.Println(result.Inspect())
}

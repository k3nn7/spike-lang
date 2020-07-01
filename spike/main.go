package main

import (
	"fmt"
	"os"
	"spike-interpreter-go/spike/eval"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
)

func main() {
	lexerInstance := lexer.New(os.Stdin)
	parserInstance := parser.New(lexerInstance)
	environment := eval.NewEnvironment()

	program, err := parserInstance.ParseProgram()
	if err != nil {
		fmt.Printf("Parser error: %s\n", err)
		return
	}

	result, err := eval.Eval(program, environment)
	if err != nil {
		fmt.Printf("Runtime error: %s\n", err)
		return
	}

	fmt.Println(result.Inspect())
}

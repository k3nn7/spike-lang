package main

import (
	"fmt"
	"os"
	"spike-interpreter-go/spike/eval"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/object"
	"spike-interpreter-go/spike/parser"
)

func main() {
	input, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Parser error: %s\n", err)
		return
	}

	lexerInstance := lexer.New(input)
	parserInstance := parser.New(lexerInstance)
	environment := object.NewEnvironment()

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

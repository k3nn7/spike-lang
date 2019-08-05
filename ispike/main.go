package main

import (
	"os"
	"spike-interpreter-go/ispike/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}

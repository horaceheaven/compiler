package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("")
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	if err := runProgram(os.Args[1]); err != nil {
		os.Exit(1)
	}

}

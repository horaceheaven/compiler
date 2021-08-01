package repl

import (
	"bufio"
	"fmt"
	"io"
)

const prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Print(prompt)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
	}
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var version = "master/unreleased"

func main() {
	eval := flag.String("eval", "", "Code to execute")
	vers := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *vers {
		fmt.Printf("monkey %s\n", version)
		os.Exit(1)
	}

	if *eval != "" {
		// TODO: add execute func
		os.Exit(1)
	}

	var input []byte
	var err error

	if len(flag.Args()) > 0 {
		input, err = ioutil.ReadFile(os.Args[1])
	} else {
		input, err = ioutil.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Printf("Error reading: %s\n", err.Error())
	}

	// TODO: add execute
	fmt.Println(input)
}

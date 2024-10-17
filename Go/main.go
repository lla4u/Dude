package main

import (
	"fmt"
	"os"

	"github.com/lla4u/Dude/cli"
)

func main() {
	// This main function is only responsible for calling the cli package
	// and handling errors returned by the cli package

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

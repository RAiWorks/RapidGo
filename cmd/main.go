package main

import (
	"fmt"
	"os"

	"github.com/RAiWorks/RapidGo/v2/core/cli"
)

func main() {
	// This is the library's built-in CLI.
	// Application projects should use cli.Set*() to wire their code.
	// See: https://github.com/RAiWorks/RapidGo-starter
	fmt.Fprintln(os.Stderr, "RapidGo is a library. Create a project with: rapidgo new myapp")
	fmt.Fprintln(os.Stderr, "Or see: https://github.com/RAiWorks/RapidGo-starter")
	cli.Execute()
}

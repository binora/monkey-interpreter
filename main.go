package main

import (
	"fmt"
	"interpreters/repl"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey Programming Language\n", u.Username)
	repl.Start(os.Stdin, os.Stdout)
}

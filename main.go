package main

import (
	"github.com/oringik/crypto-chateau/command"
	"os"
)

func main() {
	os.Exit(command.Run(os.Args[1:]))
}

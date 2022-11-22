package main

import (
	"flag"
	"github.com/oringik/crypto-chateau/gen/defs"
	"log"
)

var (
	inputFilepath  string
	outputFilepath string
	language       string
)

func init() {
	flag.StringVar(&inputFilepath, "chateau_file", "", "chateau file")
	flag.StringVar(&outputFilepath, "codegen_output", "", "codegenOutput")
	flag.StringVar(&language, "language", "", "currently supported: go, dart")
}

func main() {
	flag.Parse()

	err := defs.GenerateDefinitions(inputFilepath, outputFilepath, language)
	if err != nil {
		log.Fatalf("error generating file: " + err.Error())
	}
}

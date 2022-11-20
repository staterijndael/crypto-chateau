package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/oringik/crypto-chateau/gen/ast"
	"github.com/oringik/crypto-chateau/gen/gen"
	"github.com/oringik/crypto-chateau/gen/lexem"
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

	var (
		generator func(*ast.Ast) string
		outputExt string
	)

	switch language {
	case "":
		log.Fatal("language is not specified")
	case "go":
		generator = gen.GenerateDefinitions
		outputExt = "go"
	case "dart":
		generator = gen.GenerateDefinitionsDart
		outputExt = "dart"
	default:
		log.Fatal(language + " is not supported, only go, dart")
	}

	file, err := os.Open(inputFilepath)
	if err != nil {
		log.Fatal("input file open failed: " + err.Error())
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("input file read failed: " + err.Error())
	}

	definitions := generator(ast.GenerateAst(lexem.LexemParse(string(content))))

	outputFilename := fmt.Sprintf("%s%cgen_definitions.%s", outputFilepath, os.PathSeparator, outputExt)

	err = os.WriteFile(outputFilename, []byte(definitions), 0644)
	if err != nil {
		log.Fatal("failed to save in output file: " + err.Error())
	}

	log.Println("generated definitions saved in " + outputFilename)
}

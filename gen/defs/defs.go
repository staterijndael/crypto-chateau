package defs

import (
	"errors"
	"fmt"
	"github.com/oringik/crypto-chateau/gen/ast"
	"github.com/oringik/crypto-chateau/gen/gen"
	"github.com/oringik/crypto-chateau/gen/lexem"
	"io"
	"log"
	"os"
)

func GenerateDefinitions(inputFile string, outputFilepath string, language string) error {
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
		return errors.New(language + " is not supported, only go, dart")
	}

	file, err := os.Open(inputFile)
	if err != nil {
		return errors.New("input file open failed: " + err.Error())
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return errors.New("input file read failed: " + err.Error())
	}

	definitions := generator(ast.GenerateAst(lexem.LexemParse(string(content))))

	outputFilename := fmt.Sprintf("%s%cgen_definitions.%s", outputFilepath, os.PathSeparator, outputExt)

	err = os.WriteFile(outputFilename, []byte(definitions), 0644)
	if err != nil {
		return errors.New("failed to save in output file: " + err.Error())
	}

	log.Println("generated definitions saved in " + outputFilename)

	return nil
}

package command

import (
	"flag"
	"fmt"
	ast2 "github.com/oringik/crypto-chateau/gen/ast"
	"github.com/oringik/crypto-chateau/gen/gen"
	lexem2 "github.com/oringik/crypto-chateau/gen/lexem"
	"io"
	"os"
)

type RunOptions struct {
}

func Run(args []string) int {
	return RunCustom(args, nil)
}

func RunCustom(args []string, runOpts *RunOptions) int {
	inputFile := flag.String("chateau_file", "", "chateau file")
	outputCodegenFile := flag.String("codegen_output", "", "codegenOutput")
	language := flag.String("language", "", "language")
	flag.Parse()

	if *language == "" {
		fmt.Printf("your language flag is not set \n")

		return 2
	}
	if *language != "go" && *language != "dart" {
		fmt.Printf("language %s is unsupported; supported languages: go, dart \n", *language)

		return 2
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("cant read %s file; check your chateau_file flag \n", *inputFile)

		return 2
	}

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	lexems := lexem2.LexemParse(string(content))

	ast := ast2.GenerateAst(lexems)

	var definitionsGeneratedOutput string

	if *language == "go" {
		definitionsGeneratedOutput = gen.GenerateDefinitions(ast)
	} else if *language == "dart" {
		definitionsGeneratedOutput = gen.GenerateDefinitionsDart(ast)
	}

	if *outputCodegenFile == "" {
		*outputCodegenFile = "."
	}

	err = os.WriteFile(*outputCodegenFile+"/gen_definitions."+*language, []byte(definitionsGeneratedOutput), 0644)
	if err != nil {
		fmt.Printf("cant write your generated files; %s \n", err)

		return 2
	}

	fmt.Println("successfully generated")
	return 0
}

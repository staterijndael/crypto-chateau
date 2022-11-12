package command

import (
	"flag"
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

	if *language != "go" && *language != "dart" {
		panic("supported languages: go, dart")
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
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

	err = os.WriteFile(*outputCodegenFile+"/gen_definitions."+*language, []byte(definitionsGeneratedOutput), 0644)
	if err != nil {
		panic(err)
	}

	return 0
}

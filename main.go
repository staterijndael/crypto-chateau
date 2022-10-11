package main

import (
	"flag"
	ast2 "github.com/Oringik/crypto-chateau/gen/ast"
	"github.com/Oringik/crypto-chateau/gen/gen"
	lexem2 "github.com/Oringik/crypto-chateau/gen/lexem"
	"io/ioutil"
	"os"
)

func main() {
	inputFile := flag.String("chateau_file", "", "chateau file")
	outputCodegenFile := flag.String("codegen_output", "", "codegenOutput")
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	lexems := lexem2.LexemParse(string(content))

	ast := ast2.GenerateAst(lexems)

	definitionsGeneratedOutput := gen.GenerateDefinitions(ast)

	err = ioutil.WriteFile(*outputCodegenFile+"/gen_definitions.go", []byte(definitionsGeneratedOutput), 0644)
	if err != nil {
		panic(err)
	}
}

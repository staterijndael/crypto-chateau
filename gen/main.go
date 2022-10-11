package main

import (
	ast2 "github.com/Oringik/crypto-chateau/gen/ast"
	"github.com/Oringik/crypto-chateau/gen/gen"
	lexem2 "github.com/Oringik/crypto-chateau/gen/lexem"
	"io/ioutil"
	"os"
)

func main() {
	//inputFile := flag.String("chateauFile", "", "chateau file")
	//outputCodegenFile := flag.String("codegenOutput", "", "codegenOutput")
	//flag.Parse()

	file, err := os.Open("./gen/endpoints.chateau")
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

	//definitionsOutputFile, err := os.Create("./generated/gen_definitions.go")
	//if err != nil{
	//	panic(err)
	//}

	err = ioutil.WriteFile("./gen/generated/gen_definitions.go", []byte(definitionsGeneratedOutput), 0644)
	if err != nil {
		panic(err)
	}
}

package visitor

import (
	"fmt"
	"os"
	"wacc_32/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

const semanticError = 200

type WaccParser struct {
	*parser.WaccParser
	errorCounter *syntaxErrorCounter
}

func NewWaccParser(data string, location string) *WaccParser {
	inputStream := antlr.NewInputStream(data)
	lexer := parser.NewWaccLexer(inputStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	wp := &WaccParser{
		WaccParser: parser.NewWaccParser(tokenStream),
		errorCounter: &syntaxErrorCounter{
			antlr.NewDefaultErrorListener(),
			0,
			location,
		},
	}
	wp.AddErrorListener(wp.errorCounter)
	return wp
}

func (w *WaccParser) derive(data string, libLocation string) *WaccParser {
	inputStream := antlr.NewInputStream(data)
	lexer := parser.NewWaccLexer(inputStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	wp := &WaccParser{
		parser.NewWaccParser(tokenStream),
		w.errorCounter,
	}
	wp.AddErrorListener(wp.errorCounter)
	return wp
}

func (wp *WaccParser) GetParseTree() antlr.ParseTree {
	return wp.Program()
}

func (wp *WaccParser) PrintParseTree(tree antlr.ParseTree) {
	wp.printTreeHelper(tree, wp.GetRuleNames(), 0, "")
}

//This is based off the TreesStringTree function in antlr4/trees.go
func (wp *WaccParser) printTreeHelper(tree antlr.Tree, ruleNames []string,
	n int, indentation string) int {
	printHeader(n)
	fmt.Print(indentation)
	x := n + 1
	s := antlr.TreesGetNodeText(tree, ruleNames, nil)
	s = antlr.EscapeWhitespace(s, false)
	c := tree.GetChildCount()
	if s != "" {
		if n == 0 {
			fmt.Println(s)
		} else {
			fmt.Println("- ", s)
		}
	}
	for i := 0; i < c; i++ {
		x = wp.printTreeHelper(tree.GetChild(i), ruleNames, x, indentation+"  ")
	}
	return x
}

//Prints a line with 8 characters
func printHeader(num int) {
	s := fmt.Sprint(num)
	fmt.Print(s)
	for i := 0; i < 8-len(s); i++ {
		fmt.Print(" ")
	}
}

//SyntaxCheck exits if we have more than 1 syntax error
func (wp *WaccParser) SyntaxCheck() {
	if wp.errorCounter.errors > 0 {
		os.Exit(100)
	}
}

type syntaxErrorCounter struct {
	*antlr.DefaultErrorListener
	errors   int
	location string
}

func (sel *syntaxErrorCounter) SyntaxError(_ antlr.Recognizer, _ interface{},
	_, _ int, _ string, _ antlr.RecognitionException) {
	sel.errors++
	errMsg := "Syntax error"
	if sel.location != "" {
		errMsg += " in " + sel.location
	}
	println(errMsg)
}

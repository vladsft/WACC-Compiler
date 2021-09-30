package main

import (
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"wacc_32/assembly"
	"wacc_32/ast"
	"wacc_32/visitor"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"golang.org/x/exp/errors/fmt"
)

/************************************************************
 * THIS FILE IMPLEMENTS HELPER FUNCTIONS FOR THE FOLLOWING: *
 *  -p --parse_only                           				*
 *  -t --print_ast                           				*
 ************************************************************/

const (
	ok = iota * 100
	syntaxError
	semanticError
)

func getAST(parseTree antlr.ParseTree, waccParser *visitor.WaccParser) ast.AST {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	visitor := visitor.NewWaccVisitor("", cwd, waccParser)
	return visitor.Visit(parseTree).(ast.AST)
}

func printAST(ast ast.AST) {
	fmt.Println("===========================================================")
	fmt.Println(ast.String())
	fmt.Println("===========================================================")
}

func semanticCheck(tree ast.AST) {
	errChan := make(chan error)
	codeChan := make(chan int)
	go semanticErrorListener(errChan, codeChan)
	tree.Check(ast.Context{SemanticErrChan: errChan})
	if <-codeChan == semanticError {
		os.Exit(semanticError)
	}
}

func getLineNum(line string) int {
	start := len("Line [")
	end := strings.Index(line, ":")
	num, _ := strconv.Atoi(line[start:end])
	return num
}

func semanticErrorListener(errChan <-chan error, codeChan chan<- int) {
	code := ok
	errs := make([]string, 0)
	for err := range errChan {
		if err != nil {
			code = semanticError
			errs = append(errs, err.Error())
		}
	}
	sort.Slice(errs, func(i, j int) bool {
		return getLineNum(errs[i]) < getLineNum(errs[j])
	})
	for _, err := range errs {
		println(err)
	}
	codeChan <- code
	close(codeChan)
	switch n := len(errs); n {
	case 0:
	case 1:
		println("1 semantic error detected\n")
	default:
		fmt.Fprintf(os.Stderr, "%d semantic errors detected\n", n)
	}
}

func generateCode(tree ast.AST, writeTo string) {
	codeGen := assembly.NewArm11CodeGenerator()
	arm11Code := codeGen.GenerateCode(tree)
	err := ioutil.WriteFile(writeTo, []byte(arm11Code), 0644)
	if err != nil {
		panic(err)
	}
}

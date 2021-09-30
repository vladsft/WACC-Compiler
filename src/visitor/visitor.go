package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var _ parser.WaccParserVisitor = &WaccVisitor{}

//WaccVisitor is a BaseWaccParserVisitor with the missing methods implementes
//WaccVisitor is used to generate the AST
type WaccVisitor struct {
	*parser.BaseWaccParserVisitor
	importName string
	location   string //Used for relative addressing
	libMng     *libManager
	parser     *WaccParser
}

//NewWaccVisitor constructs a WaccVistor
func NewWaccVisitor(importName, location string, waccParser *WaccParser) *WaccVisitor {
	return &WaccVisitor{
		BaseWaccParserVisitor: &parser.BaseWaccParserVisitor{},
		importName:            importName,
		libMng:                newLibManager(),
		location:              location,
		parser:                waccParser,
	}
}

//Visit visits and returns the result of visiting the child
func (w *WaccVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(w)
}

//VisitChildren visits and returns list of children
func (w *WaccVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	var result []interface{}
	n := node.GetChildCount()
	for i := 0; i < n; i++ {
		c := node.GetChild(i).(antlr.ParseTree)
		childResult := c.Accept(w)
		result = append(result, childResult)
	}
	return result
}

//VisitNonTerminalChildren visits and returns list of all non terminal children
func (w *WaccVisitor) VisitNonTerminalChildren(node antlr.RuleNode) interface{} {
	var result []ast.AST
	n := node.GetChildCount()
	for i := 0; i < n; i++ {
		c := node.GetChild(i).(antlr.ParseTree)
		switch c.(type) {
		case antlr.TerminalNode:
			continue
		}
		childResult := c.Accept(w)
		result = append(result, childResult.(ast.AST))
	}
	return result
}

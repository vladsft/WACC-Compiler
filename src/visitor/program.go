package visitor

import (
	"io/ioutil"
	"path/filepath"
	"sync"
	"wacc_32/ast"
	"wacc_32/errors"
	"wacc_32/parser"
)

//VisitProgram returns a Program with the correct list of functions
func (w *WaccVisitor) VisitProgram(ctx *parser.ProgramContext) interface{} {
	//Visit all library functions before the main program
	go func() {
		w.VisitLibraryProgram(ctx)

		mainStatCtx := ctx.Stat()

		var mainStat ast.Statement
		var mainPos errors.Position
		if mainStatCtx != nil {
			mainStat = mainStatCtx.Accept(w).(ast.Statement)
			mainPos = getPos(ctx.Stat())
		} else {
			mainStat = ast.NewStatSkip()
			mainPos = errors.Position{}
		}

		w.libMng.functions <- ast.NewMainFunction(mainStat, mainPos)
		w.libMng.calledFunctions <- "main"
		close(w.libMng.calledFunctions)
		close(w.libMng.functions)
	}()

	funcs := make([]*ast.Function, 0) //This is a set in Go
	fNames := make(map[string]struct{})
	concurrentFNames := make(map[string]struct{})
	userTypes := make([]*ast.UserType, 0)

	for {
		select {
		case fName, ok := <-w.libMng.calledFunctions:
			if !ok {
				w.libMng.calledFunctions = nil
				break
			}
			fNames[fName] = struct{}{}
		case fName, ok := <-w.libMng.concurrentFunctions:
			if !ok {
				w.libMng.concurrentFunctions = nil
				break
			}
			concurrentFNames[fName] = struct{}{}
		case function, ok := <-w.libMng.functions:
			if !ok {
				w.libMng.functions = nil
				break
			}

			funcs = append(funcs, function)
		case userType, ok := <-w.libMng.userTypes:
			if !ok {
				w.libMng.userTypes = nil
				break
			}
			userTypes = append(userTypes, userType)
		}
		if w.libMng.functions == nil && w.libMng.calledFunctions == nil {
			break
		}
	}

	//Visit statements
	pos := getPos(ctx)

	return ast.NewProgram(userTypes, funcs, pos)
}

//Struct that stores a filepath and whether the library has been loaded
type importFilePair struct {
	filepath string
	loaded   bool
}

//VisitLibraryProgram visits all functions imports in a library file
//Doesn't visit the main function
func (w *WaccVisitor) VisitLibraryProgram(ctx *parser.ProgramContext) interface{} {
	funcsCtx := ctx.AllFunction()
	userTypesCtx := ctx.AllUserType()
	imports := ctx.AllImportfile()

	var wg sync.WaitGroup
	//Visit all imports

	for _, importFile := range imports {
		wg.Add(1)

		go func(importFile parser.IImportfileContext) {
			defer wg.Done()
			if filePair := importFile.Accept(w).(importFilePair); !filePair.loaded {
				libLocation := w.location + "/" + filePair.filepath
				data, err := ioutil.ReadFile(libLocation)
				if err != nil {
					w.throwImportError(errors.NewLibraryNotFoundError(filePair.filepath, getPos(ctx)))
				}
				wParser := w.parser.derive(string(data), libLocation)

				subVisitor := &WaccVisitor{
					BaseWaccParserVisitor: w.BaseWaccParserVisitor,
					libMng:                w.libMng.derive(),
					importName:            formatFilepath(filePair.filepath),
					location:              filepath.Dir(libLocation),
					parser:                wParser,
				}
				parseTree := wParser.GetParseTree()
				subVisitor.VisitLibraryProgram(parseTree.(*parser.ProgramContext))
			}
		}(importFile)
	}

	wg.Wait()

	//Visit all functions and send them through the channel
	for _, funcCtx := range funcsCtx {
		w.libMng.functions <- funcCtx.Accept(w).(*ast.Function)
	}
	//Visit all structs and classes and send them through the channel
	for _, userTypeCtx := range userTypesCtx {
		w.libMng.userTypes <- userTypeCtx.Accept(w).(*ast.UserType)
	}
	return nil //This won't actually be used
}

//VisitIdent returns an Ident with the correct name
func (w *WaccVisitor) VisitIdent(ctx *parser.IdentContext) interface{} {
	name := ctx.IDENT().GetText()
	return ast.NewIdent(name, getPos(ctx))
}

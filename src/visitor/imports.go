package visitor

import (
	"os"
	"regexp"
	"strings"
	"wacc_32/ast"
	"wacc_32/errors"
	"wacc_32/parser"
)

//VisitLibident returns an Ident with the correct name
func (w *WaccVisitor) VisitLibident(ctx *parser.LibidentContext) interface{} {
	fName := ctx.Fieldident().Accept(w).(*ast.Ident)
	fStr := fName.GetName()

	if ctx.ACCESSOR() != nil {
		libAlias := ctx.Ident().GetText()

		libPath, ok := w.libMng.getLib(libAlias)
		if !ok {
			w.throwImportError(errors.NewLibraryNotImportedError(libAlias, getPos(ctx)))
		}
		fStr = formatFilepath(libPath) + fStr
	} else {
		fStr = w.importName + fStr
	}
	return ast.NewIdent(fStr, getPos(ctx))
}

func (w *WaccVisitor) throwImportError(err error) {
	println(err.Error())
	os.Exit(semanticError)
}

//VisitImportfile visits and registers an imported file
func (w *WaccVisitor) VisitImportfile(ctx *parser.ImportfileContext) interface{} {
	waccFile := ctx.Waccfile().Accept(w).(waccfile)
	if ctx.AS() != nil {
		waccFile.filename = ctx.Ident().GetText()
	}
	loaded := w.libMng.addLib(waccFile.filepath, waccFile.filename)
	return importFilePair{
		filepath: waccFile.filepath,
		loaded:   loaded,
	}
}

type waccfile struct {
	filepath, filename string
}

var filepathParser = regexp.MustCompile(`\"((.*/)*(.*)\.wacc)\"`) //Index 1 to get the path and 3 to get the filename

//VisitWaccfile parses the filename of a wacc library
func (w *WaccVisitor) VisitWaccfile(ctx *parser.WaccfileContext) interface{} {
	name := ctx.STRING_LITER().GetText()
	if strings.Contains(name, "$") {
		println(errors.NewInvalidImportPathError(getPos(ctx)))
		os.Exit(100)
	}
	names := filepathParser.FindAllStringSubmatch(name, 1)[0]
	return waccfile{
		filepath: names[1],
		filename: names[3],
	}
}

package linter

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var NoExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noExitInMain",
	Doc:  "check for direct os.Exit() in main function",
	Run:  runNoExitInMain,
}

func runNoExitInMain(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		for _, decl := range file.Decls {
			fn, okFD := decl.(*ast.FuncDecl)
			if !okFD || fn.Name.Name != "main" {
				continue
			}

			ast.Inspect(fn, func(n ast.Node) bool {
				ce, okCE := n.(*ast.CallExpr)
				if !okCE {
					return true
				}

				se, okSE := ce.Fun.(*ast.SelectorExpr)
				if !okSE {
					return true
				}

				if ident, ok := se.X.(*ast.Ident); ok && ident.Name == "os" && se.Sel.Name == "Exit" {
					pass.Reportf(ce.Pos(), "avoid using os.Exit directly in main function")
				}

				return true
			})
		}
	}

	return nil, nil
}

package txndbarg

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "txndbarg",
		Doc:   "checks the dB args within database transactions",
		Run:   run,
		Flags: flagSet,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if selectorExpr.Sel.Name == "InTransaction" {
					if ident, identOK := selectorExpr.X.(*ast.Ident); identOK {
						checkInTransactionFunction(pass, f, ident.Name, callExpr.Args)
					}
				}
			}

			return true
		})
	}

	return nil, nil
}

func checkInTransactionFunction(pass *analysis.Pass, file *ast.File, dbName string, args []ast.Expr) {
	if len(args) != 2 {
		return
	}

	arg0, arg0OK := args[0].(*ast.Ident)
	if !arg0OK {
		return
	}

	field, fieldOK := arg0.Obj.Decl.(*ast.Field)
	if !fieldOK {
		return
	}

	sel, selOK := field.Type.(*ast.SelectorExpr)
	if !selOK {
		return
	}

	pkgIdent, pkgOK := sel.X.(*ast.Ident)
	if !pkgOK {
		return
	}

	if pkgIdent.Name != "context" || sel.Sel.Name != "Context" {
		return
	}

	arg1, arg1OK := args[1].(*ast.FuncLit)
	if !arg1OK {
		return
	}

	ast.Inspect(arg1, func(node ast.Node) bool {
		if ident, ok := node.(*ast.Ident); ok {
			if ident.Name == dbName {
				pass.Reportf(ident.NamePos, "cannot use %v within InTransaction function", dbName)
			}
		}

		return true
	})
}

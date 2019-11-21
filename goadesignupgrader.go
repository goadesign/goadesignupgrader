package goadesignupgrader

import (
	"go/format"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "goadesignupgrader",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "goadesignupgrader is ..."

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Remove imports for v1.
		for _, v := range []string{
			"github.com/goadesign/goa/design",
			"github.com/goadesign/goa/design/apidsl",
		} {
			astutil.DeleteImport(pass.Fset, file, v)
			astutil.DeleteNamedImport(pass.Fset, file, ".", v)
		}

		// Add imports for v3.
		astutil.AddNamedImport(pass.Fset, file, ".", "goa.design/goa/v3/dsl")

		f, err := os.OpenFile(pass.Fset.File(file.Pos()).Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if err := format.Node(f, pass.Fset, file); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

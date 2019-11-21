package goadesignupgrader

import (
	"go/ast"
	"go/format"
	"os"
	"regexp"

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

var regexpWildcard = regexp.MustCompile(`/:([a-zA-Z0-9_]+)`)

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

		astutil.Apply(file, func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.CallExpr:
				fun, ok := n.Fun.(*ast.Ident)
				if !ok {
					return true
				}
				switch fun.Name {
				case "Resource":
					// Replace Resource with Service.
					fun.Name = "Service"
				case "Action":
					// Replace Action with Method.
					fun.Name = "Method"
				case "MediaType":
					// Replace MediaType with ResultType.
					fun.Name = "ResultType"
				case "GET", " HEAD", " POST", " PUT", " DELETE", " CONNECT", " OPTIONS", " TRACE", " PATCH":
					// Replace colons with curly braces in HTTP routing DSLs.
					for _, arg := range n.Args {
						b := arg.(*ast.BasicLit)
						b.Value = replaceWildcard(b.Value)
					}
				}
			case *ast.Ident:
				switch n.Name {
				case "Integer":
					// Replace Integer with Int.
					n.Name = "Int"
				}
			}
			return true
		}, nil)

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

func replaceWildcard(s string) string {
	return regexpWildcard.ReplaceAllString(s, "/{$1}")
}

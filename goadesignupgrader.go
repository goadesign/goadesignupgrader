package goadesignupgrader

import (
	"go/format"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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

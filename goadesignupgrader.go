package goadesignupgrader

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
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
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
		(*ast.CallExpr)(nil),
		(*ast.Ident)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.ImportSpec:
			path, err := strconv.Unquote(n.Path.Value)
			if err != nil {
				return
			}
			switch path {
			case "github.com/goadesign/goa/design":
				pass.Report(analysis.Diagnostic{
					Pos: n.Pos(), Message: `"github.com/goadesign/goa/design" should be removed`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Remove", TextEdits: []analysis.TextEdit{
						{Pos: n.Pos(), End: n.End(), NewText: []byte{}},
					}}},
				})
			case "github.com/goadesign/goa/design/apidsl":
				pass.Report(analysis.Diagnostic{
					Pos: n.Path.Pos(), Message: `"github.com/goadesign/goa/design/apidsl" should be replaced with "goa.design/goa/v3/dsl"`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: n.Path.Pos(), End: n.Path.End(), NewText: []byte(`"goa.design/goa/v3/dsl"`)},
					}}},
				})
			}
		case *ast.CallExpr:
			fun, ok := n.Fun.(*ast.Ident)
			if !ok {
				return
			}
			switch fun.Name {
			case "Resource":
				pass.Report(analysis.Diagnostic{
					Pos: fun.Pos(), Message: `Resource should be replaced with Service`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: fun.Pos(), End: fun.End(), NewText: []byte("Service")},
					}}},
				})
			case "Action":
				pass.Report(analysis.Diagnostic{
					Pos: fun.Pos(), Message: `Action should be replaced with Method`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: fun.Pos(), End: fun.End(), NewText: []byte("Method")},
					}}},
				})
			case "MediaType":
				pass.Report(analysis.Diagnostic{
					Pos: fun.Pos(), Message: `MediaType should be replaced with ResultType`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: fun.Pos(), End: fun.End(), NewText: []byte("ResultType")},
					}}},
				})
			case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH":
				var (
					hasColon bool
					replaced string
					pos      token.Pos
					end      token.Pos
				)
				for _, arg := range n.Args {
					b, ok := arg.(*ast.BasicLit)
					if !ok {
						continue
					}
					replaced = replaceWildcard(b.Value)
					if replaced != b.Value {
						hasColon = true
						pos = b.Pos()
						end = b.End()
					}
				}
				if hasColon {
					pass.Report(analysis.Diagnostic{
						Pos: pos, Message: `colons in HTTP routing DSLs should be replaced with curly braces`,
						SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
							{Pos: pos, End: end, NewText: []byte(replaced)},
						}}},
					})
				}
			default:
				for _, arg := range n.Args {
					i, ok := arg.(*ast.Ident)
					if !ok {
						continue
					}
					switch i.Name {
					case "DateTime":
						i.Name = "String"
						fun, ok := n.Args[len(n.Args)-1].(*ast.FuncLit)
						if !ok {
							fun = &ast.FuncLit{
								Type: &ast.FuncType{},
								Body: &ast.BlockStmt{},
							}
							n.Args = append(n.Args, fun)
						}
						fun.Body.List = append(fun.Body.List, &ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.Ident{
									Name: "Format",
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "FormatDateTime",
									},
								},
							},
						})
						var buf bytes.Buffer
						if err := format.Node(&buf, token.NewFileSet(), n); err != nil {
							return
						}
						pass.Report(analysis.Diagnostic{
							Pos: i.Pos(), Message: `DateTime should be replaced with String + Format(FormatDateTime)`,
							SuggestedFixes: []analysis.SuggestedFix{
								{Message: "Replace", TextEdits: []analysis.TextEdit{{Pos: n.Pos(), End: n.End(), NewText: buf.Bytes()}}},
							},
						})
					}
				}
			}
		case *ast.Ident:
			switch n.Name {
			case "Integer":
				pass.Report(analysis.Diagnostic{
					Pos: n.Pos(), Message: `Integer should be replaced with Int`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: n.Pos(), End: n.End(), NewText: []byte("Int")},
					}}},
				})
			}
		}
	})

	for _, file := range pass.Files {
		astutil.Apply(file, func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.ExprStmt:
				cal, ok := n.X.(*ast.CallExpr)
				if !ok {
					return true
				}
				fun, ok := cal.Fun.(*ast.Ident)
				if !ok {
					return true
				}
				switch fun.Name {
				case "BasePath":
					// Replace BasePath with Path and move it into HTTP.
					fun.Name = "Path"
					switch nn := c.Parent().(type) {
					case *ast.BlockStmt:
						var (
							index int
							http  *ast.CallExpr
						)
						for i, v := range nn.List {
							switch nnn := v.(type) {
							case *ast.ExprStmt:
								call, ok := nnn.X.(*ast.CallExpr)
								if !ok {
									continue
								}
								funn, ok := call.Fun.(*ast.Ident)
								if !ok {
									continue
								}
								switch funn.Name {
								case "HTTP":
									http = call
								case "BasePath":
									index = i
								}
							}
						}
						if http == nil {
							http = &ast.CallExpr{
								Fun: &ast.Ident{
									Name: "HTTP",
								},
								Args: []ast.Expr{},
							}
							nn.List = append([]ast.Stmt{
								&ast.ExprStmt{
									X: http,
								},
							}, nn.List...)
							index++
						}
						var (
							ok   bool
							funn *ast.FuncLit
						)
						if len(http.Args) > 0 {
							funn, ok = http.Args[len(http.Args)-1].(*ast.FuncLit)
						}
						if !ok {
							funn = &ast.FuncLit{
								Type: &ast.FuncType{},
								Body: &ast.BlockStmt{},
							}
							http.Args = append(http.Args, funn)
						}
						funn.Body.List = append(funn.Body.List, n)
						nn.List = append(nn.List[:index], nn.List[index+1:]...)
					}
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

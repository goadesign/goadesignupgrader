package goadesignupgrader

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"regexp"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
		(*ast.BlockStmt)(nil),
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
			case "HashOf":
				pass.Report(analysis.Diagnostic{
					Pos: fun.Pos(), Message: `HashOf should be replaced with MapOf`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: fun.Pos(), End: fun.End(), NewText: []byte("MapOf")},
					}}},
				})
			case "Metadata":
				pass.Report(analysis.Diagnostic{
					Pos: fun.Pos(), Message: `Metadata should be replaced with Meta`,
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
						{Pos: fun.Pos(), End: fun.End(), NewText: []byte("Meta")},
					}}},
				})
			case "Routing":
				fun := &ast.FuncLit{
					Type: &ast.FuncType{},
					Body: &ast.BlockStmt{},
				}
				for _, arg := range n.Args {
					cal, ok := arg.(*ast.CallExpr)
					if !ok {
						continue
					}
					i, ok := cal.Fun.(*ast.Ident)
					if !ok {
						continue
					}
					switch i.Name {
					case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH":
						for _, arg := range cal.Args {
							b, ok := arg.(*ast.BasicLit)
							if !ok {
								continue
							}
							replaced := replaceWildcard(b.Value)
							if replaced != b.Value {
								b.Value = replaced
							}
						}
					}
					fun.Body.List = append(fun.Body.List, &ast.ExprStmt{
						X: cal,
					})
				}
				http := &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "HTTP",
					},
					Args: []ast.Expr{
						fun,
					},
				}
				var buf bytes.Buffer
				if err := format.Node(&buf, token.NewFileSet(), http); err != nil {
					return
				}
				pass.Report(analysis.Diagnostic{
					Pos: n.Pos(), Message: `Routing should be replaced with HTTP and colons in HTTP routing DSLs should be replaced with curly braces`,
					SuggestedFixes: []analysis.SuggestedFix{
						{Message: "Replace", TextEdits: []analysis.TextEdit{{Pos: n.Pos(), End: n.End(), NewText: buf.Bytes()}}},
					},
				})
			case "Response":
				for _, arg := range n.Args {
					i, ok := arg.(*ast.Ident)
					if !ok {
						continue
					}
					switch i.Name {
					case "Continue", "SwitchingProtocols",
						"OK", "Created", "Accepted", "NonAuthoritativeInfo", "NoContent", "ResetContent", "PartialContent",
						"MultipleChoices", "MovedPermanently", "Found", "SeeOther", "NotModified", "UseProxy", "TemporaryRedirect",
						"BadRequest", "Unauthorized", "PaymentRequired", "Forbidden", "NotFound",
						"MethodNotAllowed", "NotAcceptable", "ProxyAuthRequired", "RequestTimeout", "Conflict",
						"Gone", "LengthRequired", "PreconditionFailed", "RequestEntityTooLarge", "RequestURITooLong",
						"UnsupportedMediaType", "RequestedRangeNotSatisfiable", "ExpectationFailed", "Teapot", "UnprocessableEntity",
						"InternalServerError", "NotImplemented", "BadGateway", "ServiceUnavailable", "GatewayTimeout", "HTTPVersionNotSupported":
						name := "Status" + i.Name
						pass.Report(analysis.Diagnostic{
							Pos: i.Pos(), Message: fmt.Sprintf(`%s should be replaced with %s`, i.Name, name),
							SuggestedFixes: []analysis.SuggestedFix{
								{Message: "Replace", TextEdits: []analysis.TextEdit{{Pos: i.Pos(), End: i.End(), NewText: []byte(name)}}},
							},
						})
					}
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
		case *ast.BlockStmt:
			var (
				basePath *ast.ExprStmt
				pos      token.Pos
				end      token.Pos
			)
			for _, s := range n.List {
				stmt, ok := s.(*ast.ExprStmt)
				if !ok {
					continue
				}
				cal, ok := stmt.X.(*ast.CallExpr)
				if !ok {
					continue
				}
				fun, ok := cal.Fun.(*ast.Ident)
				if !ok {
					continue
				}
				switch fun.Name {
				case "BasePath":
					fun.Name = "Path"
					basePath = stmt
					pos = stmt.Pos()
					end = stmt.End()
				}
			}
			if basePath == nil {
				return
			}
			http := &ast.CallExpr{
				Fun: &ast.Ident{
					Name: "HTTP",
				},
				Args: []ast.Expr{
					&ast.FuncLit{
						Type: &ast.FuncType{},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{basePath},
						},
					},
				},
			}
			var buf bytes.Buffer
			if err := format.Node(&buf, token.NewFileSet(), http); err != nil {
				return
			}
			pass.Report(analysis.Diagnostic{
				Pos: pos, Message: `BasePath should be replaced with Path and move it into HTTP`,
				SuggestedFixes: []analysis.SuggestedFix{{Message: "Replace", TextEdits: []analysis.TextEdit{
					{Pos: pos, End: end, NewText: buf.Bytes()},
				}}},
			})
		}
	})

	return nil, nil
}

func replaceWildcard(s string) string {
	return regexpWildcard.ReplaceAllString(s, "/{$1}")
}

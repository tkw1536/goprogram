package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/tkw1536/pkglib/docfmt"
	"golang.org/x/tools/go/analysis"
)

const (
	goprogramPackage     = "github.com/tkw1536/goprogram"
	exitPackage          = "github.com/tkw1536/goprogram/exit"
	errorsPackage        = "errors"
	errorsNewFunc        = "New"
	descriptionType      = "Description"
	errorType            = "Error"
	messageFieldName     = "Message"
	descriptionFieldName = "Description"
	descriptionTagName   = "description"
	withMessageFunc      = "WithMessage"
	withMessageFFunc     = "WithMessageF"
)

var exceptions []string

func init() {
	values := os.Getenv("GOPROGRAM_ERRLINT_EXCEPTIONS")
	if values == "" {
		return
	}
	exceptions = strings.Split(values, ",")
}

// DocFmtAnalyzer reports incorrectly formatted calls to docfmt.
// These are:
//   - incorrectly formatted exit.Error messages and calls to WithMessage / WithMessageF
//   - incorrectly formatted ggman.Description{Description} instantiations
//   - incorrectly set description struct tags
//   - incorrect calls to errors.New
//
// Furthermore, the environment variable GOPROGRAM_ERRLINT_EXCEPTIONS may contain comma-seperated words to be ignored.
var DocFmtAnalyzer = &analysis.Analyzer{
	Name: "docfmt",
	Doc:  "reports exit.Error instances with statically unsafe messages",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			lintLiteralStructField(pass, file, exitPackage, errorType, messageFieldName, func(str string) (results []lintResult) {
				for _, err := range docfmt.Validate(str, exceptions...) {
					results = append(results, lintResult{
						Message: "message %q failed validation: %s",
						Args: []any{
							str,
							err.Error(),
						},
					})
				}
				return
			})

			lintLiteralStructField(pass, file, goprogramPackage, descriptionType, descriptionFieldName, func(str string) (results []lintResult) {
				for _, err := range docfmt.Validate(str, exceptions...) {
					results = append(results, lintResult{
						Message: "description %q failed validation: %s",
						Args: []any{
							str,
							err.Error(),
						},
					})
				}
				return
			})

			lintMethodIthStringArg(pass, file, exitPackage, errorType, withMessageFunc, 0, func(str string) (results []lintResult) {
				for _, err := range docfmt.Validate(str, exceptions...) {
					results = append(results, lintResult{
						Message: "%s(%q) failed validation: %s",
						Args: []any{
							withMessageFunc,
							str,
							err.Error(),
						},
					})
				}
				return
			})

			lintMethodIthStringArg(pass, file, exitPackage, errorType, withMessageFFunc, 0, func(str string) (results []lintResult) {
				for _, err := range docfmt.Validate(str, exceptions...) {
					results = append(results, lintResult{
						Message: "%s(%q) failed validation: %s",
						Args: []any{
							withMessageFFunc,
							str,
							err.Error(),
						},
					})
				}
				return
			})

			lintStructTag(pass, file, descriptionTagName, func(str string) (results []lintResult) {
				for _, err := range docfmt.Validate(str, exceptions...) {
					results = append(results, lintResult{
						Message: "description %q failed validation: %s",
						Args: []any{
							str,
							err.Error(),
						},
					})
				}
				return
			})

			lintFuncIthStringArg(pass, file, errorsPackage, errorsNewFunc, 0, func(str string) (results []lintResult) {
				for _, err := range docfmt.Validate(str, exceptions...) {
					results = append(results, lintResult{
						Message: "%s(%q) failed validation: %s",
						Args: []any{
							withMessageFFunc,
							str,
							err.Error(),
						},
					})
				}
				return
			})
		}

		return nil, nil
	},
}

type lintResult struct {
	Message string
	Args    []any
}

func lintStructTag(pass *analysis.Pass, file ast.Node, tag string, lint func(string) []lintResult) {
	// inspect for instantiations of ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		field, ok := n.(*ast.Field)
		if !ok {
			return true
		}

		if field.Tag == nil {
			return true
		}

		str, ok := astStringLiteral(field.Tag)
		if !ok {
			return true
		}

		val, ok := reflect.StructTag(str).Lookup(tag)
		if !ok {
			return true
		}

		// lint and report errors
		for _, res := range lint(val) {
			pass.Reportf(
				field.Tag.Pos(),
				res.Message,
				res.Args...,
			)
		}

		return true
	})
}

func lintLiteralStructField(pass *analysis.Pass, file ast.Node, pkg, tp, field string, lint func(string) []lintResult) {
	// inspect for instantiations of ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// check that we are of exit error type
		if !isStructType(pass.TypesInfo.TypeOf(comp), pkg, tp) {
			return true
		}

		var valueNode ast.Node

		// find the message key
		for _, elt := range comp.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			k, ok := kv.Key.(*ast.Ident)
			if !ok || k.Name != field {
				continue
			}

			valueNode = kv.Value
			break
		}

		// parse the message node as a basic literal
		str, ok := astStringLiteral(valueNode)
		if !ok {
			return true
		}

		// lint and report errors
		for _, res := range lint(str) {
			pass.Reportf(
				comp.Pos(),
				res.Message,
				res.Args...,
			)
		}

		return true
	})
}

func lintMethodIthStringArg(pass *analysis.Pass, file ast.Node, pkg, tp, mname string, i int, lint func(string) []lintResult) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		name, args, ok := isMethodCall(pass, call, pkg, tp)
		if !ok {
			return true
		}

		// not enough arguments
		if len(args) <= i || (name != mname) {
			return true
		}

		node := args[i]

		str, ok := astStringLiteral(node)
		if !ok {
			return true
		}

		// lint and report errors
		for _, res := range lint(str) {
			pass.Reportf(
				node.Pos(),
				res.Message,
				res.Args...,
			)
		}

		return true
	})
}

func lintFuncIthStringArg(pass *analysis.Pass, file ast.Node, pkg, fname string, i int, lint func(string) []lintResult) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		name, args, ok := isFuncCall(pass, call, pkg)
		if !ok {
			return true
		}

		if len(args) <= i || (name != fname) {
			return true
		}

		node := args[i]

		str, ok := astStringLiteral(node)
		if !ok {
			return true
		}

		// lint and report errors
		for _, res := range lint(str) {
			pass.Reportf(
				node.Pos(),
				res.Message,
				res.Args...,
			)
		}

		return true
	})
}

// isMethodCall checks if the given expression calls a method of the given (pkg, tp) struct type.
// If so returns the name of the function called, and the argument passed.
func isMethodCall(pass *analysis.Pass, call *ast.CallExpr, pkg, tp string) (string, []ast.Expr, bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", nil, false
	}

	if !isStructType(pass.TypesInfo.TypeOf(selector.X), pkg, tp) {
		return "", nil, false
	}

	return selector.Sel.Name, call.Args, true
}

// isFuncCall checks if the given expression calls a function of the given pkg.
// If so returns the name of the function called, and the argument passed.
func isFuncCall(pass *analysis.Pass, call *ast.CallExpr, pkg string) (string, []ast.Expr, bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", nil, false
	}
	identifier, ok := selector.X.(*ast.Ident)
	if !ok {
		return "", nil, false
	}
	obj := pass.TypesInfo.ObjectOf(identifier)
	if obj == nil {
		return "", nil, false
	}
	pname, ok := obj.(*types.PkgName)
	if pname == nil || !ok {
		return "", nil, false
	}

	if pname.Imported().Path() != pkg {
		return "", nil, false
	}

	return selector.Sel.Name, call.Args, true
}

// isStructType checks if tp represents a struct type in the given package with the given name
func isStructType(tp types.Type, pkg, name string) bool {
	if a := tp.Underlying(); a != tp && isStructType(a, pkg, name) {
		return true
	}

	named, ok := tp.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}
	return obj.Pkg().Path() == pkg && obj.Name() == name
}

// astStringLiteral returns the value of node as a string literal
// if the value is not a string literal, returns "", false.
func astStringLiteral(node ast.Node) (value string, ok bool) {
	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}

	// unquote the literal
	str, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", false
	}
	return str, true
}

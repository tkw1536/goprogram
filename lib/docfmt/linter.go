package docfmt

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

const (
	exitPackage      = "github.com/tkw1536/goprogram/exit"
	errorType        = "Error"
	messageFieldName = "Message"
	withMessageFunc  = "WithMessage"
	withMessageFFunc = "WithMessageF"
)

// ExitErrorAnalyzer reports incorrectly formatted exit.Error messages and calls to WithMessage / WithMessageF
var ExitErrorAnalyzer = &analysis.Analyzer{
	Name: "docfmt_error_message",
	Doc:  "reports exit.Error instances with statically unsafe messages",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {

			// inspect for instantiations of ast.Node
			ast.Inspect(file, func(n ast.Node) bool {
				comp, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}

				// check that we are of exit error type
				if !isStructType(pass.TypesInfo.TypeOf(comp), exitPackage, errorType) {
					return true
				}

				var messageNode ast.Node

				// find the message key
				for _, elt := range comp.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					k, ok := kv.Key.(*ast.Ident)
					if !ok || k.Name != messageFieldName {
						continue
					}

					messageNode = kv.Value
					break
				}

				// parse the message node as a basic literal
				str, ok := astStringLiteral(messageNode)
				if !ok {
					return true
				}

				// and validate it
				errors := Validate(str)
				if len(errors) == 0 {
					return true
				}

				for _, res := range errors {
					pass.Reportf(
						comp.Pos(),
						"message %q failed validation: %s",
						str,
						res.Error(),
					)
				}

				return true
			})

			// inspect for instantiations of ast.Node
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				name, args, ok := isCallOf(pass, call, exitPackage, errorType)
				if !ok {
					return true
				}

				if len(args) == 0 || (name != withMessageFunc && name != withMessageFFunc) {
					return true
				}

				node := args[0]

				result, ok := astStringLiteral(node)
				if !ok {
					return true
				}

				errors := Validate(result)
				for _, res := range errors {
					pass.Reportf(
						call.Pos(),
						"%s(%q,...) call failed validation: %s",
						name,
						result,
						res.Error(),
					)
				}

				return true
			})
		}

		return nil, nil
	},
}

// isCallOf checks if the given expression calls a function of the given (pkg, tp) struct type.
// If so returns the name of the function called, and the argument passed.
func isCallOf(pass *analysis.Pass, call *ast.CallExpr, pkg, tp string) (string, []ast.Expr, bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", nil, false
	}

	if !isStructType(pass.TypesInfo.TypeOf(selector.X), pkg, tp) {
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

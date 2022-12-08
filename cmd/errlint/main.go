// Command errlint lints source files with docfmt.ExitErrorAnalyzer.
package main

import (
	"github.com/tkw1536/goprogram/lib/docfmt"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(docfmt.DocFmtAnalyzer)
}

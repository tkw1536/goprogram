// Command errlint lints source files with docfmt.ExitErrorAnalyzer.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(DocFmtAnalyzer)
}

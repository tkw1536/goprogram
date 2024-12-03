// Command errlint lints source files with docfmt.ExitErrorAnalyzer.
//
//spellchecker:words main
package main

//spellchecker:words errlint docfmt

//spellchecker:words golang tools analysis singlechecker
import (
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(DocFmtAnalyzer)
}

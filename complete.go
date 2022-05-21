package goprogram

import (
	"strings"

	"github.com/tkw1536/goprogram/parser"
)

// Complete attempts to parse argv into a set of arguments to parse a specific command, and provide completion on the last argument in argv.
//
// Whenever possible, complete attempts to return a completion.
// This means that any parsing errors might be silenced.
//
// Furthermore, exact results should not be relied upon, and might be incorrect.
func (p Program[E, P, F, R]) Complete(argv []string) (pfCompletions []parser.Completion, err error) {
	// TESTME

	context := Context[E, P, F, R]{
		Program: p,
	}

	// perform completion on program flags
	hasPf, pfCompletions, pfErr := context.Args.complete(argv)

	// expand keywords and arguments
	// then attempt to load the final command
	hasKeyword, _ := context.expandKeywords()
	_, hasAlias := context.expandAliases()
	command, hasCommand := p.Command(context.Args.Command)

	switch {
	case !hasCommand && hasPf:
		// we don't have a real command, but we did get completions from the before the command
		// so we should attempt to complete those!
		return pfCompletions, pfErr
	case !hasCommand:
		// we did not get a command, and did not have completions for them
		// so we should attempt to complete the command-like argument itself!
		return p.cCommandName(context.Args.Command, !hasKeyword, !hasAlias), nil
	default:
		// we got an actual command, so we should complete the actual content
		return context.complete(command)
	}
}

// cCommandName provides completions for a command-like name
//
// A command-like query is a command, keyword or alias.
func (p Program[E, P, F, R]) cCommandName(query string, includeKeywords bool, includeAlias bool) (completions []parser.Completion) {
	for _, cmd := range p.Commands() {
		if strings.HasPrefix(cmd, query) {
			c, _ := p.Command(cmd)

			completions = append(completions, parser.Completion{
				Item:        cmd,
				Description: c.Description().Description,
			})
		}
	}

	if includeAlias {
		for _, alias := range p.Aliases() {
			if strings.HasPrefix(alias, query) {
				a := p.aliases[alias]

				completions = append(completions, parser.Completion{
					Item:        alias,
					Description: a.Description,
				})
			}
		}
	}

	if includeKeywords {
		for _, keyword := range p.Keywords() {
			if strings.HasPrefix(keyword, query) {
				completions = append(completions, parser.Completion{
					Item: keyword,
				})
			}
		}
	}

	return
}

package goprogram

import (
	"strings"

	"github.com/tkw1536/goprogram/parser"
)

// Complete provides tab completion for this program
func (p Program[E, P, F, R]) Complete(params P, argv []string) (completions []parser.Completion, err error) {
	// TESTME

	context := Context[E, P, F, R]{
		Program: p,
	}

	// perform completion on the first set of arguments!
	ok, completions, err := context.Args.completeProgramFlags(argv)

	// perform completion on the arguments
	// TODO: Share this code with Main()
	keyword, hasKeyword := p.keywords[context.Args.Command]
	if hasKeyword {
		if err := keyword(&context.Args, &context.Args.pos); err != nil {
			// FIXME(twiesing): Do we want to attempt completion anyways?
			return nil, nil
		}
	}

	// expand alias (if any)
	alias, hasAlias := p.aliases[context.Args.Command]
	if hasAlias {
		context.Args.Command, context.Args.pos = alias.Invoke(context.Args.pos)
	}

	// check to load the command
	command, hasCommand := p.Command(context.Args.Command)
	if !hasCommand {
		if ok {
			return completions, err
		}

		// user tried to invoke an unknown command
		// so we can't do any completion!
		if len(context.Args.pos) > 0 {
			return nil, nil
		}

		// else complete the list of commands
		// complete the list of commands!
		return p.completeCommandLike(context.Args.Command, !hasKeyword, !hasAlias), nil
	}

	// do the completion on the command!
	return context.complete(command)
}

func (p Program[E, P, F, R]) completeCommandLike(query string, includeKeywords bool, includeAlias bool) (completions []parser.Completion) {
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

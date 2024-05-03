//spellchecker:words goprogram
package goprogram

//spellchecker:words Positionals

// Keywords are special "commands" that manipulate arguments and positionals before execution.
//
// Keywords can not be stopped by calls to universal flags; they are expanded once before aliases and command expansion takes place.
type Keyword[F any] func(args *Arguments[F], pos *[]string) error

// RegisterKeyword registers a new keyword.
// See also Keyword.
//
// If an keyword already exists, RegisterKeyword calls panic().
func (p *Program[E, P, F, R]) RegisterKeyword(name string, keyword Keyword[F]) {
	if p.keywords == nil {
		p.keywords = make(map[string]Keyword[F])
	}

	if _, ok := p.keywords[name]; ok {
		panic("RegisterKeyword(): Keyword already registered")
	}

	p.keywords[name] = keyword
}

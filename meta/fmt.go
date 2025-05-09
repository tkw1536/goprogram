// Package meta contains facilities to provide meta-information about programs and commands.
//
//spellchecker:words meta
package meta

//spellchecker:words strings sync
import (
	"strings"
	"sync"
)

// builderPool used by various formatters in this package.
var builderPool = &sync.Pool{
	New: func() any { return new(strings.Builder) },
}

// String generates a usage page for this Meta.
func (meta Meta) String() string {
	// NOTE: This method is a wrapper around the Meta.WriteMessageTo method and is untested
	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	_ = meta.WriteMessageTo(builder) // error should never occur
	return builder.String()
}

// JoinCommands joins a list of commands into a single string.
func JoinCommands(commands []string) string {
	// NOTE: This function is untested.

	// grab a builder from the pool
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	_ = Meta{Commands: commands}.writeCommandsTo(builder) // error should never occur
	return builder.String()
}

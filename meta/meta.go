//spellchecker:words meta
package meta

//spellchecker:words strconv github pkglib docfmt
import (
	"fmt"
	"io"
	"strconv"

	"github.com/tkw1536/pkglib/docfmt"
)

//spellchecker:words positionals

// Meta holds meta-information about an entire program or a subcommand.
// It is used to generate a usage page.
type Meta struct {
	// Name of the Executable and Current command.
	// When Command is empty, the entire struct describes the program as a whole.
	Executable string
	Command    string

	// Description holds a human-readable description of the object being described.
	Description string

	// Applicable Global, Command and Positional Flags.
	GlobalFlags  []Flag
	CommandFlags []Flag
	Positionals  []Positional

	// List of available sub-commands, only set when Command == "".
	Commands []string
}

// WriteMessageTo writes the human-readable message of this meta into w.
func (meta Meta) WriteMessageTo(w io.Writer) error {
	if meta.Command != "" {
		return meta.writeCommandMessageTo(w)
	}
	return meta.writeProgramMessageTo(w)
}

// subSpec is spec for a subcommand.
const subSpec = "COMMAND [ARGS...]"

// subMsgTpl is the usage message of a subcommand.
// It consists of two parts.
const (
	// subMsgTpl = subMsg1 + "%s" + subMsg2.
	subMsg1 = "Command to call. One of "
	subMsg2 = ". See individual commands for more help."
)

func (meta Meta) writeProgramMessageTo(w io.Writer) error {
	//
	// Command specification
	//

	// main command
	if _, err := io.WriteString(w, "Usage: "); err != nil {
		return fmt.Errorf("unable to write 'Usage :': %w", err)
	}
	if _, err := io.WriteString(w, meta.Executable); err != nil {
		return fmt.Errorf("unable to write executable: %w", err)
	}

	for _, arg := range meta.GlobalFlags {
		if _, err := io.WriteString(w, " "); err != nil {
			return fmt.Errorf("unable to write ' ': %w", err)
		}
		if err := arg.WriteSpecTo(w); err != nil {
			return fmt.Errorf("unable to write spec: %w", err)
		}
	}

	if _, err := io.WriteString(w, " [--] "); err != nil {
		return fmt.Errorf("unable to write ' [--] ': %w", err)
	}
	if _, err := io.WriteString(w, subSpec); err != nil {
		return fmt.Errorf("unable to write sub spec: %w", err)
	}

	// description (if any)
	if meta.Description != "" {
		if _, err := io.WriteString(w, "\n\n"); err != nil {
			return fmt.Errorf("unable to write newlines: %w", err)
		}
		if _, err := io.WriteString(w, docfmt.Format(meta.Description)); err != nil {
			return fmt.Errorf("unable to write description: %w", err)
		}
	}

	//
	// Argument description
	//

	for _, arg := range meta.GlobalFlags {
		if err := arg.WriteMessageTo(w); err != nil {
			return err
		}
	}

	// write a usage message for the commands

	if _, err := io.WriteString(w, usageMsg1); err != nil {
		return fmt.Errorf("unable to write usage message: %w", err)
	}
	if _, err := io.WriteString(w, subSpec); err != nil {
		return fmt.Errorf("unable to write sub specification: %w", err)
	}
	if _, err := io.WriteString(w, usageMsg2); err != nil {
		return fmt.Errorf("unable to write usage message: %w", err)
	}

	// replace the list of commands in subMsgTpl
	if _, err := io.WriteString(w, subMsg1); err != nil {
		return fmt.Errorf("unable to write sub specification: %w", err)
	}
	if err := meta.writeCommandsTo(w); err != nil {
		return fmt.Errorf("unable to write commands: %w", err)
	}
	if _, err := io.WriteString(w, subMsg2); err != nil {
		return fmt.Errorf("unable to sub specification: %w", err)
	}

	if _, err := io.WriteString(w, usageMsg3); err != nil {
		return fmt.Errorf("unable to write usage message: %w", err)
	}

	return nil
}

// WriteCommandsTo writes the list of commands to w.
func (meta Meta) writeCommandsTo(w io.Writer) error {
	if len(meta.Commands) == 0 {
		return nil
	}
	if _, err := io.WriteString(w, strconv.Quote(meta.Commands[0])); err != nil {
		return fmt.Errorf("unable to write quoted command: %w", err)
	}
	for _, cmd := range meta.Commands[1:] {
		if _, err := io.WriteString(w, ", "); err != nil {
			return fmt.Errorf("unable to write ' ': %w", err)
		}
		if _, err := io.WriteString(w, strconv.Quote(cmd)); err != nil {
			return fmt.Errorf("unable to write quoted command: %w", err)
		}
	}
	return nil
}

func (page Meta) writeCommandMessageTo(w io.Writer) error {
	//
	// Command specification
	//

	// main command
	if _, err := io.WriteString(w, "Usage: "); err != nil {
		return fmt.Errorf("unable to write 'Usage :': %w", err)
	}
	if _, err := io.WriteString(w, page.Executable); err != nil {
		return fmt.Errorf("unable to write executable: %w", err)
	}

	for _, arg := range page.GlobalFlags {
		if _, err := io.WriteString(w, " "); err != nil {
			return fmt.Errorf("unable to write ' ': %w", err)
		}
		if err := arg.WriteSpecTo(w); err != nil {
			return fmt.Errorf("unable to write argument spec: %w", err)
		}
	}

	if len(page.GlobalFlags) > 0 {
		if _, err := io.WriteString(w, " [--]"); err != nil {
			return fmt.Errorf("unable to write ' [--]': %w", err)
		}
	}

	// subcommand
	if _, err := io.WriteString(w, " "); err != nil {
		return fmt.Errorf("unable to write ' ': %w", err)
	}
	if _, err := io.WriteString(w, page.Command); err != nil {
		return fmt.Errorf("unable to write command: %w", err)
	}

	for _, arg := range page.CommandFlags {
		if _, err := io.WriteString(w, " "); err != nil {
			return fmt.Errorf("unable to write ' ': %w", err)
		}
		if err := arg.WriteSpecTo(w); err != nil {
			return fmt.Errorf("unable to write flag spec: %w", err)
		}
	}

	if len(page.Positionals) != 0 {
		if _, err := io.WriteString(w, " [--]"); err != nil {
			return fmt.Errorf("unable to write ' [--]': %w", err)
		}

		for _, p := range page.Positionals {
			if _, err := io.WriteString(w, " "); err != nil {
				return fmt.Errorf("unable to write ' ': %w", err)
			}
			if err := p.WriteSpecTo(w); err != nil {
				return fmt.Errorf("unable to write positional spec: %w", err)
			}
		}
	}

	// description (if any)
	if page.Description != "" {
		if _, err := io.WriteString(w, "\n\n"); err != nil {
			return fmt.Errorf("unable to write newlines: %w", err)
		}
		if _, err := io.WriteString(w, docfmt.Format(page.Description)); err != nil {
			return fmt.Errorf("unable to write description: %w", err)
		}
	}

	//
	// Argument description
	//

	if _, err := io.WriteString(w, "\n\nGlobal Arguments:"); err != nil {
		return fmt.Errorf("unable to write 'Global Arguments': %w", err)
	}
	for _, opt := range page.GlobalFlags {
		if err := opt.WriteMessageTo(w); err != nil {
			return fmt.Errorf("unable to write global flags: %w", err)
		}
	}

	// no command arguments provided!
	if len(page.CommandFlags) == 0 && len(page.Positionals) == 0 {
		return nil
	}

	if _, err := io.WriteString(w, "\n\nCommand Arguments:"); err != nil {
		return fmt.Errorf("unable to write 'Command Arguments': %w", err)
	}

	for _, opt := range page.CommandFlags {
		if err := opt.WriteMessageTo(w); err != nil {
			return fmt.Errorf("unable to write command flags: %w", err)
		}
	}

	for _, p := range page.Positionals {
		if _, err := io.WriteString(w, usageMsg1); err != nil {
			return fmt.Errorf("unable to write positional usage message: %w", err)
		}
		if err := p.WriteSpecTo(w); err != nil {
			return fmt.Errorf("unable to write positional: %w", err)
		}
		if _, err := io.WriteString(w, usageMsg2); err != nil {
			return fmt.Errorf("unable to write positional usage message: %w", err)
		}
		if _, err := io.WriteString(w, docfmt.Format(p.Usage)); err != nil {
			return fmt.Errorf("unable to write positional usage message: %w", err)
		}
		if _, err := io.WriteString(w, usageMsg3); err != nil {
			return fmt.Errorf("unable to write positional usage message: %w", err)
		}
	}
	return nil
}

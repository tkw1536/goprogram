// Package docfmt implements formatting and checking of user format strings.
//
// Strings are checked at runtime for proper formatting
// Checking is disabled by default, but can be enabled by building with the "doccheck" tag.
// See Check.
package docfmt

import (
	"fmt"
	"strings"
)

// Check checks the message for proper formatting.
//
// When checking is disabled, no checking is performed.
// When checking is enabled and a message fails to pass validation, calls panic()
func Check(message string) {
	if enabled {
		if errors := Validate(message); len(errors) != 0 {
			panic(&CheckError{
				Message: message,
				Results: errors,
			})
		}
	}
}

// CheckError is returned when a message fails validation.
// It implements the built-in error interface.
type CheckError struct {
	Results []ValidationResult

	// message is the message being checked
	Message string
}

func (ce CheckError) Error() string {
	// NOTE(twiesing): This function is untested because it is used only for developing

	messages := make([]string, len(ce.Results))
	for i, res := range ce.Results {
		messages[i] = res.Error()
	}

	return fmt.Sprintf("message %q failed validation: %s", ce.Message, strings.Join(messages, "\n"))
}

package docfmt

import (
	"fmt"
	"strings"
)

// AssertValid asserts that message is propertly format and calling Validate on it returns no results.
//
// When checking is disabled, no runtime checking is performed.
// When checking is enabled and a message fails to pass validation, calls panic()
func AssertValid(message string) {
	if enabled {
		if errors := Validate(message); len(errors) != 0 {
			panic(&ValidationError{
				Message: message,
				Results: errors,
			})
		}
	}
}

// ValidationError is returned when a message fails validation.
// It implements the built-in error interface.
type ValidationError struct {
	Results []ValidationResult

	// message is the message being checked
	Message string
}

func (ce ValidationError) Error() string {
	// NOTE(twiesing): This function is untested because it is used only for developing

	messages := make([]string, len(ce.Results))
	for i, res := range ce.Results {
		messages[i] = res.Error()
	}

	return fmt.Sprintf("message %q failed validation: %s", ce.Message, strings.Join(messages, "\n"))
}

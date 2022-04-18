// Package doccheck checks documentation strings at runtime for proper formatting.
//
// Checking is disabled by default, but can be enabled by building with the "doccheck" tag.
package doccheck

import (
	"fmt"
	"strings"
	"unicode"
)

// Check checks the message for proper formatting.
//
// When checking is disabled, no checking is performed.
// When checking is enabled and a message fails to pass validation, calls panic()
func Check(message string) {
	if enabled {
		if err := Validate(message); err != nil {
			panic(err)
		}
	}
}

// CheckError is returned when a message fails validation.
// It implements the built-in error interface.
type CheckError struct {
	// message is the message being checked
	Message string

	// Part and Index give information about the part in which the failure occured
	Part  string
	Index int

	// Failure is the failure that occured
	Failure string
}

func (ce CheckError) Error() string {
	// NOTE(twiesing): This function is untested because it is used only for developing
	return fmt.Sprintf("message %q failed validation: part %q: %s", ce.Message, ce.Part, ce.Failure)
}

// Validate validates that message is formatted correctly.
// When the message passes validation returns nil, otherwise a *CheckError.
//
// To perform validation, the message is first split into parts delimited by ':'s.
//
// Then the following tests are performed for each part:
//   - a part must be non-empty
//   - a part must start with a space, unless it is the first part
//   - a part may not contain extra spaces at the start and end
//   - a part may not end with a period
//   - a part may not start with a capital letter, unless all letters of the first word are capital
//
// An empty message is excluded from all checks and passes without error.
func Validate(message string) error {
	if message == "" { // zero message is ok
		return nil
	}
	for i, part := range strings.Split(message, ":") {
		failure := validatePart(i, part)
		if failure != "" {
			return &CheckError{
				Message: message,

				Part:  part,
				Index: i,

				Failure: failure,
			}
		}
	}
	return nil
}

func validatePart(i int, part string) string {
	if len(part) == 0 {
		return "empty"
	}
	if i != 0 {
		if part[0] != ' ' {
			return "missing a leading space"
		}
		part = part[1:]
		if len(part) == 0 {
			return "empty"
		}
	}

	if strings.TrimSpace(part) != part {
		return "contains extra space"
	}
	if strings.HasSuffix(part, ".") {
		return "ends with a period"
	}
	if unicode.IsUpper(rune(part[0])) {
		word := strings.Fields(part)[0]
		allUpper := strings.IndexFunc(word, func(r rune) bool {
			return !unicode.IsUpper(r)
		}) == -1
		if !allUpper {
			return "starts with upper case"
		}
	}
	return ""
}

//go:build !doccheck

package doccheck

// Enabled checks if the doccheck package is enabled.
//
// It will return false because it is currently disabled.
func Enabled() bool {
	return false
}

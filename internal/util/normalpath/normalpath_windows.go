//go:build windows
// +build windows

package normalpath

import (
	"path/filepath"
	"strings"
)

// NormalizeAndValidate normalizes and validates the given path.
//
// This calls Normalize on the path.
// Returns Error if the path is not relative or jumps context.
// This can be used to validate that paths are valid to use with Buckets.
// The error message is safe to pass to users.
func NormalizeAndValidate(path string) (string, error) {
	normalizedPath := Normalize(path)
	if filepath.IsAbs(normalizedPath) || (len(normalizedPath) > 0 && normalizedPath[0] == '/') {
		// the stdlib implementation of `IsAbs` assumes that a volume name is required for a path to
		// be absolute, however Windows treats a `/` (normalized) rooted path as absolute _within_ the current volume.
		// In the context of validating that a path is _not_ relative, we need to reject a path that begins
		// with `/`.
		return "", NewError(path, errNotRelative)
	}
	// https://github.com/bufbuild/buf/issues/51
	if strings.HasPrefix(normalizedPath, normalizedRelPathJumpContextPrefix) {
		return "", NewError(path, errOutsideContextDir)
	}
	return normalizedPath, nil
}

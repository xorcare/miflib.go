package norm

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

// String it's function provides a normalized representation of a string.
func String(s string) string {
	old := s
	s = strings.ToValidUTF8(s, "")
	s = strings.TrimSpace(s)
	s = norm.NFKC.String(s)

	if s == old {
		return s
	}

	return String(s)
}

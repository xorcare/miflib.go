// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"bytes"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// String it's function provides a normalized representation of a string.
func String(s string) string {
	old := s
	s = strings.ToValidUTF8(s, "")
	s = norm.NFKC.String(s)
	s = rewriteNotPrintable(s)
	s = strings.ReplaceAll(s, "  ", " ")
	s = strings.TrimSpace(s)

	if s == old {
		return s
	}

	return String(s)
}

func rewriteNotPrintable(s string) string {
	buf := bytes.NewBuffer(make([]byte, 0, len(s)))
	for _, r := range s {
		switch {
		case !strconv.IsPrint(r):
			// skip not printable rune.
		default:
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jstring

import (
	"github.com/xorcare/miflib.go/internal/norm"
)

// String is the normalized string type.
type String string

// String implements the fmt.Stringer interface.
func (s String) String() string {
	return norm.String(string(s))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s String) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *String) UnmarshalText(text []byte) error {
	*s = String(norm.String(string(text)))
	return nil
}

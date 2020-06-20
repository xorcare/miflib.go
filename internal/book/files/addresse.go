// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

import (
	"github.com/xorcare/miflib.go/internal/jstring"
)

// Address contains the parameters of the file by which you can download it.
type Address struct {
	URL      string         `json:"url,omitempty"`
	Size     uint           `json:"size,omitempty"`
	Duration string         `json:"duration,omitempty"`
	Title    jstring.String `json:"title,omitempty"`
}

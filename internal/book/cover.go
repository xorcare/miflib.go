// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package book

// Cover address on the cover of the book.
type Cover struct {
	Small string `json:"small,omitempty"`
	Large string `json:"large,omitempty"`
}

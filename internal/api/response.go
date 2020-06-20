// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/xorcare/miflib.go/internal/book"
)

// ListResponse response from API which contains books and their number.
type ListResponse struct {
	Books []book.Book `json:"books,omitempty"`
	Total uint
}

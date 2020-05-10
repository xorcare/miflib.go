// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package books

import (
	"github.com/xorcare/miflib.go/internal/books/book"
)

// Books response from API which contains books and their number.
type Books struct {
	Books []book.Book `json:"books,omitempty"`
	Total uint
}

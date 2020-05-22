package api

import (
	"github.com/xorcare/miflib.go/internal/book"
)

// ListResponse response from API which contains books and their number.
type ListResponse struct {
	Books []book.Book `json:"books,omitempty"`
	Total uint
}

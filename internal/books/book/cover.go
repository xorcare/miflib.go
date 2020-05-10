package book

// Cover address on the cover of the book.
type Cover struct {
	Small string `json:"small,omitempty"`
	Large string `json:"large,omitempty"`
}

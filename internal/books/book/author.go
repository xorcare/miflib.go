package book

// Author contains information about the author of the book.
type Author struct {
	Name  string `json:"name,omitempty"`
	Photo string `json:"photo,omitempty"`
	Info  string `json:"info,omitempty"`
}

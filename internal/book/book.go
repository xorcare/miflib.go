package book

import (
	"github.com/xorcare/miflib.go/internal/book/files"
)

// Book the dataset about the book.
type Book struct {
	ID           int             `json:"id"`
	Title        string          `json:"title"`
	Subtitle     string          `json:"subtitle"`
	BookPartLink string          `json:"bookPartLink"`
	Badges       []string        `json:"badges,omitempty"`
	SimilarBooks []string        `json:"similarBooks,omitempty"`
	Cover        Cover           `json:"cover"`
	NewCover     string          `json:"newCover"`
	Files        files.Files     `json:"files"`
	TopSmile     interface{}     `json:"topSmile,omitempty"`
	MifURL       string          `json:"mifUrl"`
	Description  string          `json:"description"`
	Stickers     []TTS           `json:"stickers,omitempty"`
	Quotes       []TTS           `json:"quotes,omitempty"`
	Experts      []interface{}   `json:"experts,omitempty"`
	Photos       []files.Address `json:"photos,omitempty"`
	Videos       []files.Address `json:"videos,omitempty"`
	Spreads      []interface{}   `json:"spreads,omitempty"`
	Authors      []Author        `json:"authors,omitempty"`
	DiscountURL  string          `json:"discountUrl"`
	Downloads    int             `json:"downloads"`
}

// TTS structure for different descriptions.
type TTS struct {
	Style string `json:"style,omitempty"`
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
}

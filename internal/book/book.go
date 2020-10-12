// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package book

import (
	"github.com/xorcare/miflib.go/internal/jstring"
)

// Book the dataset about the book.
type Book struct {
	ID int `json:"id"`

	Title       jstring.String `json:"title"`
	Subtitle    string         `json:"subtitle"`
	Description string         `json:"description"`

	Authors      []Author      `json:"authors,omitempty"`
	Badges       []string      `json:"badges,omitempty"`
	BookPartLink string        `json:"bookPartLink"`
	Cover        Cover         `json:"cover"`
	DiscountURL  string        `json:"discountUrl"`
	Downloads    int           `json:"downloads"`
	Experts      []interface{} `json:"experts,omitempty"`
	Files        Files         `json:"files"`
	MifURL       string        `json:"mifUrl"`
	NewCover     string        `json:"newCover"`
	Photos       []Address     `json:"photos,omitempty"`
	Quotes       []TTS         `json:"quotes,omitempty"`
	SimilarBooks []string      `json:"similarBooks,omitempty"`
	Spreads      []interface{} `json:"spreads,omitempty"`
	Stickers     []TTS         `json:"stickers,omitempty"`
	TopSmile     interface{}   `json:"topSmile,omitempty"`
	Videos       []Address     `json:"videos,omitempty"`
}

// TTS structure for different descriptions.
type TTS struct {
	Style string         `json:"style,omitempty"`
	Text  string         `json:"text,omitempty"`
	Title jstring.String `json:"title,omitempty"`
}

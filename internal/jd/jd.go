// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jd

import (
	"encoding/json"
	"fmt"
)

var _ json.Unmarshaler = (*Addresses)(nil)

// Address contains the parameters of the file by which you can download it.
type Address struct {
	URL      string `json:"url,omitempty"`
	Size     uint   `json:"size,omitempty"`
	Duration string `json:"duration,omitempty"`
	Title    string `json:"title,omitempty"`
}

// Addresses just a set of addresses with the ability to parse not a set,
// but one address as such cases are in the API.
type Addresses []Address

// UnmarshalJSON implements json.Unmarshaler.
func (a *Addresses) UnmarshalJSON(bs []byte) error {
	if a == nil {
		return fmt.Errorf("%T: UnmarshalJSON on nil pointer", a)
	}

	address := Address{}
	if err := json.Unmarshal(bs, &address); err == nil {
		*a = append((*a)[:], address)
		return nil
	}

	addresses := make([]Address, 0, 1)
	err := json.Unmarshal(bs, &addresses)
	*a = append((*a)[:], Addresses(addresses)...)

	return err
}

// Author contains information about the author of the book.
type Author struct {
	Name  string `json:"name,omitempty"`
	Photo string `json:"photo,omitempty"`
	Info  string `json:"info,omitempty"`
}

// Book the dataset about the book.
type Book struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Subtitle     string        `json:"subtitle"`
	BookPartLink string        `json:"bookPartLink"`
	Badges       []string      `json:"badges,omitempty"`
	SimilarBooks []string      `json:"similarBooks,omitempty"`
	Cover        Cover         `json:"cover"`
	NewCover     string        `json:"newCover"`
	Files        Files         `json:"files"`
	TopSmile     interface{}   `json:"topSmile,omitempty"`
	MifURL       string        `json:"mifUrl"`
	Description  string        `json:"description"`
	Stickers     []TTS         `json:"stickers,omitempty"`
	Quotes       []TTS         `json:"quotes,omitempty"`
	Experts      []interface{} `json:"experts,omitempty"`
	Photos       []Address     `json:"photos,omitempty"`
	Videos       []Address     `json:"videos,omitempty"`
	Spreads      []interface{} `json:"spreads,omitempty"`
	Authors      []Author      `json:"authors,omitempty"`
	DiscountURL  string        `json:"discountUrl"`
	Downloads    int           `json:"downloads"`
}

// Books response from API which contains books and their number.
type Books struct {
	Books []Book `json:"books,omitempty"`
	Total uint
}

// Cover address on the cover of the book.
type Cover struct {
	Small string `json:"small,omitempty"`
	Large string `json:"large,omitempty"`
}

// Files information about all available files for download.
type Files struct {
	Books      map[string]Addresses `json:"ebook"`
	AudioBooks map[string]Addresses `json:"audiobook"`
	Demo       map[string]Addresses `json:"demo"`
}

// TTS structure for different descriptions.
type TTS struct {
	Style string `json:"style,omitempty"`
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
}

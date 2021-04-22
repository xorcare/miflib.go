// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package book

// Files information about all available files for download.
type Files struct {
	Books      Formats `json:"ebook"`
	AudioBooks Formats `json:"audiobook"`
	Demo       Formats `json:"demo"`
}

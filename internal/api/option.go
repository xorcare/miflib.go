// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

// Option it's a interface for options func.
type Option func(*Client)

// OptDoer it's option for set http client.
func OptDoer(d doer) Option {
	return func(client *Client) {
		client.http = d
	}
}

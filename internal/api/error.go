// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Error contains an error response from the server.
type Error struct {
	// Code is the HTTP response status code and will always be populated.
	Code int
	// Body is the raw response returned by the server.
	// It is often but not always JSON, depending on how the request fails.
	Body string

	URL url.URL
}

func (e *Error) Error() string {
	return fmt.Sprintf("downloader: got HTTP response of url %s code %d with body: %v", e.URL.String(), e.Code, e.Body)
}

// checkResponse returns an error (of type *Error) if the response
// status code is not 2xx.
func checkResponse(response *http.Response) error {
	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}
	body, _ := ioutil.ReadAll(response.Body)
	err := &Error{
		Code: response.StatusCode,
		Body: string(body),
		URL:  *response.Request.URL,
	}
	err.URL.User = nil
	return err
}

// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flags

import (
	"strings"
)

// Flag names constants.
const (
	Username                  = "username"
	Password                  = "password"
	Hostname                  = "hostname"
	Directory                 = "directory"
	NumThreads                = "num-threads"
	HTTPResponseHeaderTimeout = "http-response-header-timeout"
	HTTPTimeout               = "http-timeout"
	Verbose                   = "verbose"
)

// Env it's a function for conversion flag name to env variable name.
func Env(flag string, envs ...string) []string {
	return append(
		[]string{"MIFLIB_" + strings.ToUpper(strings.ReplaceAll(flag, "-", "_"))},
		envs...,
	)
}

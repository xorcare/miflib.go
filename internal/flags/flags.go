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
)

// Env it's a function for conversion flag name to env variable name.
func Env(flag string, envs ...string) []string {
	return append(
		[]string{"MIFLIB_" + strings.ToUpper(strings.ReplaceAll(flag, "-", "_"))},
		envs...,
	)
}

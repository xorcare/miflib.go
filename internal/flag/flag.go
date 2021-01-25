package flag

import (
	"runtime"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/xorcare/miflib.go/internal/flags"
)

func init() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "print help",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "print the version",
	}
}

// Username is a instance of cli flag.
var Username = &cli.StringFlag{
	Name:     flags.Username,
	Aliases:  []string{"u"},
	Usage:    "username for the library",
	Required: true,
	EnvVars:  flags.Env(flags.Username),
}

// Password is a instance of cli flag.
var Password = &cli.StringFlag{
	Name:     flags.Password,
	Aliases:  []string{"p"},
	Usage:    "password for the library",
	Required: true,
	EnvVars:  flags.Env(flags.Password),
}

// Hostname is a instance of cli flag.
var Hostname = &cli.StringFlag{
	Name:     flags.Hostname,
	Aliases:  []string{"h"},
	Usage:    "hostname for the library",
	Required: true,
	EnvVars:  flags.Env(flags.Hostname),
}

// Directory is a instance of cli flag.
var Directory = &cli.StringFlag{
	Name:    flags.Directory,
	Aliases: []string{"d"},
	Usage:   "the directory where books will be placed",
	EnvVars: flags.Env(flags.Directory),
	Value:   ".",
}

// NumThreads is a instance of cli flag.
var NumThreads = &cli.IntFlag{
	Name:    flags.NumThreads,
	Aliases: []string{"n"},
	Usage:   "number of books processed in parallel",
	EnvVars: flags.Env(flags.NumThreads),
	Value:   runtime.NumCPU(),
}

// HTTPResponseHeaderTimeout is a instance of cli flag.
var HTTPResponseHeaderTimeout = &cli.DurationFlag{
	Name: flags.HTTPResponseHeaderTimeout,
	Usage: "specifies the amount of time to wait for a server's" +
		" response headers after fully writing the request (including" +
		" its body, if any). This time does not include the time to" +
		" read the response body.",
	EnvVars: flags.Env(flags.HTTPResponseHeaderTimeout),
	Value:   time.Minute,
}

// HTTPTimeout is a instance of cli flag.
var HTTPTimeout = &cli.DurationFlag{
	Name:    flags.HTTPTimeout,
	Usage:   "timeout specifies a time limit for requests made by this tool.",
	EnvVars: flags.Env(flags.HTTPTimeout),
	Value:   time.Hour,
}

// Verbose is a instance of cli flag.
var Verbose = &cli.BoolFlag{
	Name:    flags.Verbose,
	Aliases: []string{"v"},
	EnvVars: flags.Env(flags.Verbose),
	Value:   false,
}

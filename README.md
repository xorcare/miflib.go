# miflib.go

[![travis-ci.org](https://travis-ci.org/xorcare/miflib.go.svg?branch=master)][TCI]
[![codecov.io](https://codecov.io/gh/xorcare/miflib.go/badge.svg)][COV]
[![goreportcard.com](https://goreportcard.com/badge/github.com/xorcare/miflib.go)][GRC]
[![godoc.org](https://godoc.org/github.com/xorcare/miflib.go?status.svg)][DOC]

Application to download data from [miflib](https://www.mann-ivanov-ferber.ru/b2b/elibrary) library.

## Installation

```bash
go get github.com/xorcare/miflib.go/cmd/miflib
```

## Command line interface, [CLI]

To get help on working with the program, run the command:

```bash
miflib --help
```

Example result of command execution:

```text
NAME:
   miflib - Application to download data from miflib library.

USAGE:
   miflib [global options] command [command options] [arguments...]

VERSION:
   v0.1.0-30-g15fcae6

AUTHOR:
   Vasiliy Vasilyuk <xorcare@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --username value, -u value            username for the library [$MIFLIB_USERNAME]
   --password value, -p value            password for the library [$MIFLIB_PASSWORD]
   --hostname value, -h value            hostname for the library [$MIFLIB_HOSTNAME]
   --directory value, -d value           the directory where books will be placed (default: ".") [$MIFLIB_DIRECTORY]
   --num-threads value, -n value         number of books processed in parallel (default: 12) [$MIFLIB_NUM_THREADS]
   --http-response-header-timeout value  specifies the amount of time to wait for a server's response headers after fully writing the request (including its body, if any). This time does not include the time to read the response body. (default: 1m0s) [$MIFLIB_HTTP_RESPONSE_HEADER_TIMEOUT]
   --http-timeout value                  timeout specifies a time limit for requests made by this tool. (default: 1h0m0s) [$MIFLIB_HTTP_TIMEOUT]
   --verbose, -v                         (default: false) [$MIFLIB_VERBOSE]
   --help                                print help (default: false)
   --version                             print the version (default: false)

COPYRIGHT:
   Copyright (c) 2019-2020 Vasiliy Vasilyuk
```

## License

Released under the [BSD 3-Clause License][LIC].

[LIC]: https://github.com/xorcare/miflib.go/blob/master/LICENSE 'BSD 3-Clause "New" or "Revised" License'
[TCI]: https://travis-ci.org/xorcare/miflib.go 'Travis CI is a hosted continuous integration service used to build and test software projects hosted at GitHub'
[COV]: https://codecov.io/gh/xorcare/miflib.go/branch/master 'Codecov is a code coverage tool, which is available for GitHub, Bitbucket and GitLab'
[GRC]: https://goreportcard.com/report/github.com/xorcare/miflib.go 'A web application that generates a report on the quality of an open source go project'
[DOC]: https://godoc.org/github.com/xorcare/miflib.go 'GoDoc hosts documentation for Go packages on Bitbucket, GitHub, Google Project Hosting and Launchpad'
[CLI]: https://en.wikipedia.org/wiki/Command-line_interface 'Command-line interface'

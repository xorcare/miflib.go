// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"log"
	"os"

	"github.com/xorcare/miflib.go/internal/cli"
)

// Version of the application is installed from outside during assembly.
var Version = "unknown"

func main() {
	if err := cli.New(Version).Run(os.Args); err != nil {
		log.SetFlags(0)
		log.Fatalf("error: %s", err)
	}
}

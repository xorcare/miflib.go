// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/urfave/cli/v2"
	"golang.org/x/net/publicsuffix"

	"github.com/xorcare/miflib.go/internal/downloader"
	"github.com/xorcare/miflib.go/internal/jd"
)

// Version of the application is installed from outside during assembly.
var Version = "v0.0.0"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	http.DefaultClient.Jar = jar

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "print help",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "print the version",
	}
}

func main() {
	app := &cli.App{
		Name:    "miflib",
		Action:  action,
		Version: Version,
		Authors: []*cli.Author{{
			Name:  "Vasiliy Vasilyuk",
			Email: "xorcare@gmail.com",
		}},
	}

	app.Copyright = "Copyright (c) 2019 Vasiliy Vasilyuk. All rights reserved.\n"
	app.Usage = "Application to download data from miflib library."
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "u",
			Aliases:  []string{"username"},
			Usage:    "username for the library",
			Required: true,
			EnvVars:  []string{"MIFLIB_USERNAME"},
		},
		&cli.StringFlag{
			Name:     "p",
			Aliases:  []string{"password"},
			Usage:    "password for the library",
			Required: true,
			EnvVars:  []string{"MIFLIB_PASSWORD"},
		},
		&cli.StringFlag{
			Name:     "h",
			Aliases:  []string{"hostname"},
			Usage:    "hostname for the library",
			Required: true,
			EnvVars:  []string{"MIFLIB_HOSTNAME"},
		},
		&cli.IntFlag{
			Name:    "n",
			Aliases: []string{"num-threads"},
			Usage:   "number of books processed in parallel",
			EnvVars: []string{"MIFLIB_NUM_THREADS"},
			Value:   runtime.NumCPU(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	uri, err := url.Parse(fmt.Sprintf("https://%s/auth/login.ajax", c.String("h")))
	if err != nil {
		return err
	}

	res, err := http.Post(uri.String(), "application/json;charset=utf-8", bytes.NewBufferString(
		fmt.Sprintf(`{"email":%q,"password":%q}`, c.String("u"), c.String("p")),
	))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := downloader.CheckResponse(res); err != nil {
		return err
	}

	uri.Path = "books/list.ajax"
	res, err = http.Get(uri.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := downloader.CheckResponse(res); err != nil {
		return err
	}

	books := jd.Books{}
	err = json.NewDecoder(res.Body).Decode(&books)
	if err != nil {
		return err
	}

	ch := make(chan jd.Book)
	wg := sync.WaitGroup{}
	wg.Add(c.Int("n"))
	for i := 0; i < c.Int("n"); i++ {
		go worker(&wg, ch, "books")
	}

	for _, book := range books.Books {
		ch <- book
	}

	close(ch)
	wg.Wait()

	log.Println("correct completion of processing")

	return nil
}

func worker(wg *sync.WaitGroup, books <-chan jd.Book, basepath string) {
	defer wg.Done()

	for book := range books {
		log.Println("start processing book:", book.Title, book.ID)
		func(payload interface{}) {
			defer log.Println("finish processing book:", book.Title, book.ID)
			book := payload.(jd.Book)

			basepath = path.Join(basepath, fmt.Sprintf("%05d %s", book.ID, book.Title))
			if err := os.MkdirAll(basepath, 0755); err != nil {
				log.Fatal(err)
			}

			filepath := path.Join(basepath, "book.json")

			if _, err := os.Stat(filepath); !os.IsNotExist(err) {
				log.Println("book is already downloaded earlier:", book.Title, book.ID)
				return
			}

			if err := downloader.Download(basepath, book); err != nil {
				log.Fatal(err)
			}
			file, err := os.Create(filepath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			if encoder.Encode(book) != nil {
				log.Fatal(err)
			}

			return
		}(book)

		log.Println("the book is loaded:", book.Title, book.ID)
	}
}

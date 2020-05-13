// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/signal"
	"path"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"

	"github.com/xorcare/miflib.go/internal/books"
	"github.com/xorcare/miflib.go/internal/books/book"
	"github.com/xorcare/miflib.go/internal/downloader"
)

// Version of the application is installed from outside during assembly.
var Version = "v0.0.0"

func init() {
	log.SetFlags(log.LstdFlags)

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
		&cli.StringFlag{
			Name:    "d",
			Aliases: []string{"directory"},
			Usage:   "the directory where books will be placed",
			EnvVars: []string{"MIFLIB_DIRECTORY"},
			Value:   ".",
		},
		&cli.IntFlag{
			Name:    "n",
			Aliases: []string{"num-threads"},
			Usage:   "number of books processed in parallel",
			EnvVars: []string{"MIFLIB_NUM_THREADS"},
			Value:   runtime.NumCPU(),
		},
		&cli.DurationFlag{
			Name: "http-response-header-timeout",
			Usage: "specifies the amount of time to wait for a server's" +
				" response headers after fully writing the request (including" +
				" its body, if any). This time does not include the time to" +
				" read the response body.",
			EnvVars: []string{"MIFLIB_HTTP_RESPONSE_HEADER_TIMEOUT"},
			Value:   time.Minute,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}
	http.DefaultClient.Jar = jar
	http.DefaultClient.Transport = &http.Transport{
		ResponseHeaderTimeout: c.Duration("http-response-header-timeout"),
	}

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

	ch := make(chan book.Book)

	ctx, done := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("miflib is shutting down by os interrupt signal...")
		done()
		// You need to completely subtract the channel for successful completion
		// in the event of an interruption of the program.
		for range ch {
		}
	}()

	wg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < c.Int("n"); i++ {
		wg.Go(func() error {
			return worker(ctx, ch, c.String("d"))
		})
	}

	wg.Go(func() error {
		defer close(ch)
		uri.Path = "books/list.ajax"
		res, err = http.Get(uri.String())
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if err := downloader.CheckResponse(res); err != nil {
			return err
		}

		bks := books.Books{}
		err = json.NewDecoder(res.Body).Decode(&bks)
		if err != nil {
			return err
		}

		for len(bks.Books) > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				ch <- bks.Books[0]
				bks.Books = bks.Books[1:]
			}
		}

		return nil
	})

	if err := wg.Wait(); err != nil {
		return err
	}

	log.Println("correct completion of processing")

	return nil
}

func worker(ctx context.Context, ch <-chan book.Book, basepath string) (err error) {
	for bk := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Println("start processing book:", bk.Title, bk.ID)
			err = func() error {
				defer log.Println("finish processing book:", bk.Title, bk.ID)

				basepath = path.Join(basepath, fmt.Sprintf("%05d %s", bk.ID, bk.Title))
				if err := os.MkdirAll(basepath, 0755); err != nil {
					return err
				}

				filepath := path.Join(basepath, "book.json")

				if _, err := os.Stat(filepath); !os.IsNotExist(err) {
					log.Println("book is already downloaded earlier:", bk.Title, bk.ID)
					return nil
				}

				if err := downloader.Download(ctx, basepath, bk); err != nil {
					return err
				}
				file, err := os.Create(filepath)
				if err != nil {
					return err
				}
				defer file.Close()

				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "\t")
				if encoder.Encode(bk) != nil {
					return err
				}

				return nil
			}()
		}
		if err != nil {
			return err
		}

		log.Println("the book is loaded:", bk.Title, bk.ID)
	}

	return nil
}

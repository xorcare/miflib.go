package cli

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

// New returns new instance of miflib application.
func New(version string) *cli.App {
	app := &cli.App{
		Name:    "miflib",
		Action:  action,
		Version: version,
		Authors: []*cli.Author{{
			Name:  "Vasiliy Vasilyuk",
			Email: "xorcare@gmail.com",
		}},
	}

	app.Copyright = "Copyright (c) 2019-2020 Vasiliy Vasilyuk\n"
	app.Usage = "Application to download data from miflib library."
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     flags.Username,
			Aliases:  []string{"u"},
			Usage:    "username for the library",
			Required: true,
			EnvVars:  flags.Env(flags.Username),
		},
		&cli.StringFlag{
			Name:     flags.Password,
			Aliases:  []string{"p"},
			Usage:    "password for the library",
			Required: true,
			EnvVars:  flags.Env(flags.Password),
		},
		&cli.StringFlag{
			Name:     flags.Hostname,
			Aliases:  []string{"h"},
			Usage:    "hostname for the library",
			Required: true,
			EnvVars:  flags.Env(flags.Hostname),
		},
		&cli.StringFlag{
			Name:    flags.Directory,
			Aliases: []string{"d"},
			Usage:   "the directory where books will be placed",
			EnvVars: flags.Env(flags.Directory),
			Value:   ".",
		},
		&cli.IntFlag{
			Name:    flags.NumThreads,
			Aliases: []string{"n"},
			Usage:   "number of books processed in parallel",
			EnvVars: flags.Env(flags.NumThreads),
			Value:   runtime.NumCPU(),
		},
		&cli.DurationFlag{
			Name: flags.HTTPResponseHeaderTimeout,
			Usage: "specifies the amount of time to wait for a server's" +
				" response headers after fully writing the request (including" +
				" its body, if any). This time does not include the time to" +
				" read the response body.",
			EnvVars: flags.Env(flags.HTTPResponseHeaderTimeout),
			Value:   time.Second * 10,
		},
	}

	return app
}

func action(c *cli.Context) error {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}
	http.DefaultClient.Jar = jar
	http.DefaultClient.Transport = &http.Transport{
		ResponseHeaderTimeout: c.Duration(flags.HTTPResponseHeaderTimeout),
	}

	uri, err := url.Parse(fmt.Sprintf("https://%s/auth/login.ajax", c.String(flags.Hostname)))
	if err != nil {
		return err
	}

	res, err := http.Post(uri.String(), "application/json;charset=utf-8", bytes.NewBufferString(
		fmt.Sprintf(`{"email":%q,"password":%q}`, c.String(flags.Username), c.String(flags.Password)),
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

	for i := 0; i < c.Int(flags.NumThreads); i++ {
		wg.Go(func() error {
			return worker(ctx, ch, c.String(flags.Directory))
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

				bookpath := path.Join(basepath, fmt.Sprintf("%05d %s", bk.ID, bk.Title))
				if err := os.MkdirAll(bookpath, 0755); err != nil {
					return err
				}

				filepath := path.Join(bookpath, "book.json")

				if _, err := os.Stat(filepath); !os.IsNotExist(err) {
					log.Println("book is already downloaded earlier:", bk.Title, bk.ID)
					return nil
				}

				if err := downloader.Download(ctx, bookpath, bk); err != nil {
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

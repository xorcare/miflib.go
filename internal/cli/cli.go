package cli

import (
	"context"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"

	"github.com/xorcare/miflib.go/internal/api"
	"github.com/xorcare/miflib.go/internal/book"
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

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return err
	}

	apiClient := api.NewClient(
		"https://"+c.String(flags.Hostname),
		api.OptDoer(&http.Client{
			Timeout: time.Hour,
			Transport: &http.Transport{
				ResponseHeaderTimeout: c.Duration(flags.HTTPResponseHeaderTimeout),
			},
			Jar: jar,
		}),
	)

	if err := apiClient.Login(ctx, c.String(flags.Username), c.String(flags.Password)); err != nil {
		return err
	}

	wg, ctx := errgroup.WithContext(ctx)

	loader := downloader.NewLoader(c.String(flags.Directory), apiClient)
	for i := 0; i < c.Int(flags.NumThreads); i++ {
		wg.Go(func() error {
			return loader.Worker(ctx, ch)
		})
	}

	wg.Go(func() error {
		defer close(ch)
		bks, err := apiClient.List(ctx)
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

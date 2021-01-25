package cli

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"sort"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"

	"github.com/xorcare/miflib.go/internal/api"
	"github.com/xorcare/miflib.go/internal/book"
	"github.com/xorcare/miflib.go/internal/downloader"
	"github.com/xorcare/miflib.go/internal/flag"
)

// New returns new instance of miflib application.
func New(version string) *cli.App {
	app := &cli.App{
		Name:    "miflib",
		Action:  action,
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "Vasiliy Vasilyuk",
				Email: "xorcare@gmail.com",
			},
		},
	}

	app.Copyright = "Copyright (c) 2019-2020 Vasiliy Vasilyuk\n"
	app.Usage = "Application to download data from miflib library."

	app.Flags = []cli.Flag{
		flag.Username,
		flag.Password,
		flag.Hostname,
		flag.Directory,
		flag.NumThreads,
		flag.HTTPResponseHeaderTimeout,
		flag.HTTPTimeout,
		flag.Verbose,
	}

	return app
}

func action(c *cli.Context) error {
	loggerConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          "console",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	}

	if c.Bool(flag.Verbose.Name) {
		loggerConf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, _ := loggerConf.Build()
	sugar := logger.Sugar()
	defer logger.Sync()

	ch := make(chan book.Book)

	ctx, done := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Info("miflib is shutting down by os interrupt signal...")
		done()
		// You need to completely subtract the channel for successful completion
		// in the event of an interruption of the program.
		for range ch {
		}
	}()
	jar, err := cookiejar.New(
		&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		},
	)
	if err != nil {
		return err
	}

	apiClient := api.NewClient(
		"https://"+c.String(flag.Hostname.Name),
		sugar,
		api.OptDoer(
			&http.Client{
				Timeout: c.Duration(flag.HTTPTimeout.Name),
				Transport: &http.Transport{
					ResponseHeaderTimeout: c.Duration(flag.HTTPResponseHeaderTimeout.Name),
				},
				Jar: jar,
			},
		),
	)

	if err := apiClient.Login(
		ctx,
		c.String(flag.Username.Name),
		c.String(flag.Password.Name),
	); err != nil {
		return err
	}

	wg, ctx := errgroup.WithContext(ctx)

	loader := downloader.NewLoader(c.String(flag.Directory.Name), apiClient, sugar)
	for i := 0; i < c.Int(flag.NumThreads.Name); i++ {
		wg.Go(
			func() error {
				return loader.Worker(ctx, ch)
			},
		)
	}

	wg.Go(
		func() error {
			defer close(ch)
			bks, err := apiClient.List(ctx)
			if err != nil {
				return err
			}

			sort.Slice(
				bks.Books, func(i, j int) bool {
					return bks.Books[i].ID < bks.Books[j].ID
				},
			)

			sugar.Infof("currently %d books are available for download", bks.Total)

			for i, bk := range bks.Books {
				sugar.Infof("%d books are waiting to be downloaded", int(bks.Total)-i)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- bk:
				}
			}

			return nil
		},
	)

	if err := wg.Wait(); err != nil {
		return err
	}

	logger.Info("correct completion of downloading")

	return nil
}

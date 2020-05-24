// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/xorcare/miflib.go/internal/api"
	"github.com/xorcare/miflib.go/internal/book"
	"github.com/xorcare/miflib.go/internal/book/files"
)

type logger interface {
	Infof(msg string, keysAndValues ...interface{})
	Debugf(msg string, keysAndValues ...interface{})
	Warnf(msg string, keysAndValues ...interface{})
}

// Downloader this is the file loader interface.
type Downloader interface {
	DownloadFile(ctx context.Context, url, filename string) (err error)
}

// Loader is an implementation of a handler for loading all possible materials
// from a book.
type Loader struct {
	api  Downloader
	root string
	log  logger
}

// NewLoader creates new instance of loader.
func NewLoader(basepath string, downloader Downloader, logger logger) Loader {
	return Loader{
		api:  downloader,
		root: basepath,
		log:  logger,
	}
}

// download starting the download mechanism.
func (l *Loader) download(ctx context.Context, basepath string, bk book.Book) error {
	type downloader func(context.Context, string, book.Book) error
	var downloaders = []downloader{
		l.downloadAudiobook,
		l.downloadBook,
		l.downloadCover,
		l.downloadDemo,
		l.downloadPhotos,
	}

	for _, f := range downloaders {
		if err := f(ctx, basepath, bk); err != nil {
			if err, ok := err.(*url.Error); ok {
				if err.Err.Error() == "stopped after 10 redirects" {
					l.log.Warnf("skip redirect error for the book %q", bk.Title)
					continue
				}
			}

			if err, ok := err.(*api.Error); ok && err.Code == 404 {
				l.log.Warnf("skip undiscovered files for the book %q with error %q", bk.Title, err)
				continue
			}

			return err
		}
	}

	return nil
}

// Worker it's a method for processing a channel with books,
// it downloads information for all books read from the channel.
func (l *Loader) Worker(ctx context.Context, ch <-chan book.Book) (err error) {
	for bk := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			l.log.Infof("start downloading the book %q", bk.Title)

			bookpath := path.Join(l.root, fmt.Sprintf("%05d %s", bk.ID, bk.Title))
			if err := os.MkdirAll(bookpath, 0755); err != nil {
				return err
			}

			filepath := path.Join(bookpath, "book.json")

			if _, err := os.Stat(filepath); !os.IsNotExist(err) {
				l.log.Infof("the book %q is already downloaded earlier", bk.Title)
				continue
			}

			if err := l.download(ctx, bookpath, bk); err != nil {
				return err
			}

			l.log.Infof("finishing downloading the book: %q", bk.Title)

			file, err := os.Create(filepath)
			if err != nil {
				return err
			}
			_ = file

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			if err := encoder.Encode(bk); err != nil {
				file.Close()
				return err
			}
			file.Close()

			l.log.Infof("the book %q is loaded", bk.Title)
		}
	}

	return nil
}

func (l *Loader) downloadAudiobook(ctx context.Context, basepath string, book book.Book) error {
	l.log.Infof("start downloading are audiobook for the book %q, ", book.Title)
	l.log.Debugf("available audiobook %s", book.Files.AudioBooks)
	defer l.log.Infof("finishing downloading are audiobook for the book %q, ", book.Title)
	basepath = path.Join(basepath, "audiobook")
	for key, as := range book.Files.AudioBooks {
		// The zip file contains all mp3 recordings together so there is
		// no need to download everything together.
		if key == "mp3" && len(book.Files.AudioBooks["zip"]) > 0 {
			l.log.Infof("skip mp3 because zip exists for the book %q", book.Title)
			continue
		}
		for _, address := range as {
			if err := l.downloadByAddress(ctx, basepath, key, address, book); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Loader) downloadBook(ctx context.Context, basepath string, book book.Book) error {
	l.log.Infof("start downloading are ebook for the book %q, ", book.Title)
	l.log.Debugf("available ebook %s", book.Files.Books)
	defer l.log.Infof("finishing downloading are ebook for the book %q, ", book.Title)
	basepath = path.Join(basepath, "e-book")
	for key, as := range book.Files.Books {
		for _, address := range as {
			if err := l.downloadByAddress(ctx, basepath, key, address, book); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Loader) downloadCover(ctx context.Context, basepath string, book book.Book) error {
	l.log.Infof("start downloading are cover for the book %q, ", book.Title)
	defer l.log.Infof("finishing downloading are cover for the book %q, ", book.Title)
	if err := l.downloadFileByURL(ctx, book.Cover.Large, basepath); err != nil {
		return err
	}

	return l.downloadFileByURL(ctx, book.Cover.Small, basepath)
}

func (l *Loader) downloadDemo(ctx context.Context, basepath string, book book.Book) error {
	l.log.Infof("start downloading are demo for the book %q", book.Title)
	l.log.Debugf("available demo %s", book.Files.Demo)
	defer l.log.Infof("finishing downloading are demo for the book %q, ", book.Title)
	basepath = path.Join(basepath, "demo")
	for key, as := range book.Files.Demo {
		for _, address := range as {
			if err := l.downloadByAddress(ctx, basepath, key, address, book); err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *Loader) downloadPhotos(ctx context.Context, basepath string, book book.Book) error {
	l.log.Infof("start downloading are photos for the book %q, ", book.Title)
	defer l.log.Infof("finishing downloading are photos for the book %q, ", book.Title)
	basepath = path.Join(basepath, "photos")
	for _, as := range book.Photos {
		if err := l.downloadFileByURL(ctx, as.URL, basepath); err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) downloadByAddress(ctx context.Context, basepath, ext string, ad files.Address, book book.Book) error {
	title := ad.Title
	if title == "" {
		title = book.Title
	}
	msg := fmt.Sprintf("%s.%s", title, ext)

	return l.api.DownloadFile(ctx, ad.URL, path.Join(basepath, ext, msg))
}

func (l *Loader) downloadFileByURL(ctx context.Context, url, basepath string) error {
	return l.api.DownloadFile(ctx, url, path.Join(basepath, path.Base(url)))
}

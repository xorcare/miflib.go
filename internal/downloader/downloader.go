// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/xorcare/miflib.go/internal/api"
	"github.com/xorcare/miflib.go/internal/book"
	"github.com/xorcare/miflib.go/internal/book/files"
)

// Downloader this is the file loader interface.
type Downloader interface {
	DownloadFile(ctx context.Context, url, filename string) (err error)
}

// Loader is an implementation of a handler for loading all possible materials
// from a book.
type Loader struct {
	api  Downloader
	root string
}

// NewLoader creates new instance of loader.
func NewLoader(basepath string, downloader Downloader) Loader {
	return Loader{
		api:  downloader,
		root: basepath,
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
			if er, ok := err.(*url.Error); ok {
				if er.Err.Error() == "stopped after 10 redirects" {
					log.Println("skip redirect error", err)
					continue
				}
			}

			if err, ok := err.(*api.Error); ok && err.Code == 404 {
				log.Println("skip undiscovered files", err)
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
			log.Println("start processing book:", bk.Title, bk.ID)
			err = func() error {
				defer log.Println("finish processing book:", bk.Title, bk.ID)

				bookpath := path.Join(l.root, fmt.Sprintf("%05d %s", bk.ID, bk.Title))
				if err := os.MkdirAll(bookpath, 0755); err != nil {
					return err
				}

				filepath := path.Join(bookpath, "book.json")

				if _, err := os.Stat(filepath); !os.IsNotExist(err) {
					log.Println("book is already downloaded earlier:", bk.Title, bk.ID)
					return nil
				}

				if err := l.download(ctx, bookpath, bk); err != nil {
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

func (l *Loader) downloadAudiobook(ctx context.Context, basepath string, book book.Book) error {
	log.Println("start download audiobook", book.ID)
	defer log.Println("finish download audiobook", book.ID)
	basepath = path.Join(basepath, "audiobook")
	for key, as := range book.Files.AudioBooks {
		// The zip file contains all mp3 recordings together so there is
		// no need to download everything together.
		if key == "mp3" && len(book.Files.AudioBooks["zip"]) > 0 {
			log.Println("skip mp3 if zip exists")
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
	log.Println("start download e-book", book.ID)
	defer log.Println("finish download e-book", book.ID)
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
	log.Println("start download cover", book.ID)
	defer log.Println("finish download cover", book.ID)
	if err := l.downloadFileByURL(ctx, book.Cover.Large, basepath); err != nil {
		return err
	}

	return l.downloadFileByURL(ctx, book.Cover.Small, basepath)
}

func (l *Loader) downloadDemo(ctx context.Context, basepath string, book book.Book) error {
	log.Println("start download demo", book.ID)
	defer log.Println("finish download demo", book.ID)
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
	log.Println("start download photos", book.ID)
	defer log.Println("finish download photos", book.ID)
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

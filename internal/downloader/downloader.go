// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/xorcare/miflib.go/internal/books/book"
	"github.com/xorcare/miflib.go/internal/books/book/files"
)

func init() {
	http.DefaultClient.CheckRedirect = CheckRedirect
	http.DefaultClient.Timeout = time.Hour
}

// Downloader type for handlers.
type Downloader func(ctx context.Context, path string, book book.Book) error

// DownloadHandlers list of handlers used.
var DownloadHandlers = []Downloader{
	downloadCover,
	downloadAudiobook,
	downloadBook,
	downloadDemo,
	downloadPhotos,
}

var errStoppedAfterRedirects = errors.New("stopped after 10 redirects")

// CheckRedirect custom cycle redirect error.
func CheckRedirect(_ *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errStoppedAfterRedirects
	}
	return nil
}

// Download starting the download mechanism.
func Download(ctx context.Context, basepath string, book book.Book) error {
	for _, f := range DownloadHandlers {
		if err := f(ctx, basepath, book); err != nil {
			if er, ok := err.(*url.Error); ok {
				if er.Err == errStoppedAfterRedirects {
					log.Println("skip redirect error", err)
					continue
				}
			}

			if strings.Contains(err.Error(), "got HTTP response code 404 with body") {
				log.Println("skip undiscovered files", err)
				continue
			}

			return err
		}
	}

	return nil
}

func downloadAudiobook(ctx context.Context, basepath string, book book.Book) error {
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
			if err := downloadByAddress(ctx, basepath, key, address, book); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadBook(ctx context.Context, basepath string, book book.Book) error {
	log.Println("start download e-book", book.ID)
	defer log.Println("finish download e-book", book.ID)
	basepath = path.Join(basepath, "e-book")
	for key, as := range book.Files.Books {
		for _, address := range as {
			if err := downloadByAddress(ctx, basepath, key, address, book); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadCover(ctx context.Context, basepath string, book book.Book) error {
	log.Println("start download cover", book.ID)
	defer log.Println("finish download cover", book.ID)
	if err := downloadFileByURL(ctx, book.Cover.Large, basepath); err != nil {
		return err
	}

	return downloadFileByURL(ctx, book.Cover.Small, basepath)
}

func downloadDemo(ctx context.Context, basepath string, book book.Book) error {
	log.Println("start download demo", book.ID)
	defer log.Println("finish download demo", book.ID)
	basepath = path.Join(basepath, "demo")
	for key, as := range book.Files.Demo {
		for _, address := range as {
			if err := downloadByAddress(ctx, basepath, key, address, book); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadPhotos(ctx context.Context, basepath string, book book.Book) error {
	log.Println("start download photos", book.ID)
	defer log.Println("finish download photos", book.ID)
	basepath = path.Join(basepath, "photos")
	for _, as := range book.Photos {
		if err := downloadFileByURL(ctx, as.URL, basepath); err != nil {
			return err
		}
	}

	return nil
}

func downloadByAddress(ctx context.Context, basepath, ext string, ad files.Address, book book.Book) error {
	title := ad.Title
	if title == "" {
		title = book.Title
	}
	msg := fmt.Sprintf("%s.%s", title, ext)

	return downloadFile(ctx, ad.URL, path.Join(basepath, ext, msg))
}

func downloadFile(ctx context.Context, url, filename string) error {
	log.Println("start download from url:", url, "to file:", filename)
	defer log.Println("finish download from url:", url, "to file:", filename)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	var done func()
	ctx, done = context.WithCancel(ctx)
	defer done()

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := CheckResponse(res); err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)

	return err
}

func downloadFileByURL(ctx context.Context, url, basepath string) error {
	return downloadFile(ctx, url, path.Join(basepath, path.Base(url)))
}

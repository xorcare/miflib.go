// Copyright (c) 2020 Vasiliy Vasilyuk All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package downloader

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xorcare/golden"
	"go.uber.org/zap"

	"github.com/xorcare/miflib.go/internal/api"
	"github.com/xorcare/miflib.go/internal/book"
	"github.com/xorcare/miflib.go/internal/ctxtest"
)

type apiMock struct {
	mock.Mock
}

func (d *apiMock) DownloadFile(ctx context.Context, url, filename string) (err error) {
	return d.Called(ctx, url, filename).Error(0)
}

func TestLoader_downloadAudiobook(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)
	l := &Loader{
		api:  amk,
		root: "",
		log:  zap.NewNop().Sugar(),
	}

	ctx := context.Background()

	amk.On("DownloadFile", ctx, "https://zip", "jedi/audiobook/zip/Джедайские техники.zip").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://m4b", "jedi/audiobook/m4b/Джедайские техники.m4b").Return(nil).Once()

	require.NoError(t, l.downloadAudiobook(ctx, "jedi", book.Book{
		Title: "Джедайские техники\n\r\t!",
		Files: book.Files{
			AudioBooks: map[string]book.Addresses{
				"m4b": {
					book.Address{
						URL: "https://m4b",
					},
				},
				"mp3": {
					book.Address{
						URL:   "https://mp3/0",
						Title: "Введение",
					}, book.Address{
						URL:   "https://mp3/1",
						Title: "Глава 1",
					}, book.Address{
						URL:   "https://mp3/2",
						Title: "Приложения",
					},
				},
				"ogg": {
					book.Address{
						URL:   "https://ogg/0",
						Title: "Введение\r",
					}, book.Address{
						URL:   "https://ogg/1",
						Title: "Глава 1?",
					}, book.Address{
						URL:   "https://ogg/2",
						Title: "Приложения!",
					},
				},
				"zip": {
					book.Address{
						URL: "https://zip",
					},
				},
			},
		},
	}))

	t.Run("error", func(t *testing.T) {
		amk.On("DownloadFile", ctx, "https://zip/error", "jedi/audiobook/zip/Джедайские техники.zip").Return(io.EOF).Once()
		require.Error(t, l.downloadAudiobook(ctx, "jedi", book.Book{
			Title: "Джедайские техники",
			Files: book.Files{
				AudioBooks: map[string]book.Addresses{
					"zip": {
						book.Address{
							URL: "https://zip/error",
						},
					},
				},
			},
		}))
	})
}

func TestLoader_downloadBook(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)
	l := &Loader{
		api:  amk,
		root: "",
		log:  zap.NewNop().Sugar(),
	}

	ctx := context.Background()

	amk.On("DownloadFile", ctx, "https://epub", "jedi/e-book/epub/Джедайские техники.epub").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://fb2", "jedi/e-book/fb2/Джедайские техники.fb2").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://mobi", "jedi/e-book/mobi/Джедайские техники.mobi").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://pdf", "jedi/e-book/pdf/Джедайские техники.pdf").Return(nil).Once()

	require.NoError(t, l.downloadBook(ctx, "jedi", book.Book{
		Title: "Джедайские техники",
		Files: book.Files{
			Books: map[string]book.Addresses{
				"epub": {
					book.Address{
						URL: "https://epub",
					},
				},
				"fb2": {
					book.Address{
						URL: "https://fb2",
					},
				},
				"mobi": {
					book.Address{
						URL: "https://mobi",
					},
				},
				"pdf": {
					book.Address{
						URL: "https://pdf",
					},
				},
			},
		},
	}))

	t.Run("error", func(t *testing.T) {
		amk.On("DownloadFile", ctx, "https://epub/error", "jedi/e-book/epub/Джедайские техники.epub").Return(io.EOF).Once()
		require.Error(t, l.downloadBook(ctx, "jedi", book.Book{
			Title: "Джедайские техники",
			Files: book.Files{
				Books: map[string]book.Addresses{
					"epub": {
						book.Address{
							URL: "https://epub/error",
						},
					},
				},
			},
		}))
	})
}

func TestLoader_downloadCover(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)
	l := &Loader{
		api:  amk,
		root: "",
		log:  zap.NewNop().Sugar(),
	}

	ctx := context.Background()

	amk.On("DownloadFile", ctx, "https://big.png", "jedi/big.png").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://s.png", "jedi/s.png").Return(nil).Once()

	require.NoError(t, l.downloadCover(ctx, "jedi", book.Book{
		Title: "Джедайские техники",
		Cover: book.Cover{
			Small: "https://s.png",
			Large: "https://big.png",
		},
	}))

	t.Run("error", func(t *testing.T) {
		amk.On("DownloadFile", ctx, "https://big.png", "jedi/big.png").Return(io.EOF).Once()
		require.Error(t, l.downloadCover(ctx, "jedi", book.Book{
			Title: "Джедайские техники",
			Cover: book.Cover{
				Small: "https://s.png",
				Large: "https://big.png",
			},
		}))
	})
}

func TestLoader_downloadDemo(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)
	l := &Loader{
		api:  amk,
		root: "",
		log:  zap.NewNop().Sugar(),
	}

	ctx := context.Background()

	amk.On("DownloadFile", ctx, "https://epub", "jedi/demo/epub/Джедайские техники.epub").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://fb2", "jedi/demo/fb2/Джедайские техники.fb2").Return(nil).Once()

	require.NoError(t, l.downloadDemo(ctx, "jedi", book.Book{
		Title: "Джедайские техники",
		Files: book.Files{
			Demo: map[string]book.Addresses{
				"epub": {
					book.Address{
						URL: "https://epub",
					},
				},
				"fb2": {
					book.Address{
						URL: "https://fb2",
					},
				},
			},
		},
	}))

	t.Run("error", func(t *testing.T) {
		amk.On("DownloadFile", ctx, "https://epub/error", "jedi/demo/epub/Джедайские техники.epub").Return(io.EOF).Once()
		require.Error(t, l.downloadDemo(ctx, "jedi", book.Book{
			Title: "Джедайские техники",
			Files: book.Files{
				Demo: map[string]book.Addresses{
					"epub": {
						book.Address{
							URL: "https://epub/error",
						},
					},
				},
			},
		}))
	})
}

func TestLoader_downloadPhotos(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)
	l := &Loader{
		api:  amk,
		root: "",
		log:  zap.NewNop().Sugar(),
	}

	ctx := context.Background()

	amk.On("DownloadFile", ctx, "https://032dt.png", "jedi/photos/032dt.png").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://035dt.png", "jedi/photos/035dt.png").Return(nil).Once()

	require.NoError(t, l.downloadPhotos(ctx, "jedi", book.Book{
		Title: "Джедайские техники",
		Photos: []book.Address{
			{
				URL: "https://032dt.png",
			}, {
				URL: "https://035dt.png",
			},
		},
	}))

	t.Run("error", func(t *testing.T) {
		amk.On("DownloadFile", ctx, "https://032dt.png", "jedi/photos/032dt.png").Return(io.EOF).Once()
		require.Error(t, l.downloadPhotos(ctx, "jedi", book.Book{
			Title: "Джедайские техники",
			Photos: []book.Address{
				{
					URL: "https://032dt.png",
				},
			},
		}))
	})
}

func TestLoader_download(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)

	l := &Loader{
		api:  amk,
		root: "",
		log:  zap.NewNop().Sugar(),
	}

	amk.On("DownloadFile", ctxtest.Match, "https://cover/small.png", "jedi/small.png").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://cover/large.png", "jedi/large.png").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://zip", "jedi/audiobook/zip/Джедайские техники.zip").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://pdf", "jedi/e-book/pdf/Джедайские техники.pdf").Return(nil).Once()

	amk.On("DownloadFile", ctxtest.Match, "https://fb2", "jedi/e-book/fb2/Джедайские техники.fb2").
		Return(&api.Error{Code: 404}).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://mobi", "jedi/e-book/mobi/Джедайские техники.mobi").
		Return(&url.Error{Err: errors.New("stopped after 10 redirects")}).Once()

	require.NoError(t, l.download(ctxtest.Background(), "jedi", book.Book{
		Title: "Джедайские техники",
		Photos: []book.Address{{
			URL: "https://photos/photos.png",
		}},
		Cover: book.Cover{
			Small: "https://cover/small.png",
			Large: "https://cover/large.png",
		},
		Files: book.Files{
			Books: map[string]book.Addresses{
				"pdf": {
					book.Address{
						URL: "https://pdf",
					},
				},
				"mobi": {
					book.Address{
						URL: "https://mobi",
					},
				},
				"fb2": {
					book.Address{
						URL: "https://fb2",
					},
				},
			},
			AudioBooks: map[string]book.Addresses{
				"zip": {
					book.Address{
						URL: "https://zip",
					},
				},
			},
			Demo: map[string]book.Addresses{
				"epub": {
					book.Address{
						URL: "https://epub",
					},
				},
			},
		},
	}))
}

func TestLoader_Worker(t *testing.T) {
	amk := new(apiMock)
	amk.Test(t)
	defer amk.AssertExpectations(t)

	tempDir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)

	l := &Loader{
		api:  amk,
		root: tempDir,
		log:  zap.NewNop().Sugar(),
	}

	amk.On("DownloadFile", ctxtest.Match, "https://cover/small.png", filepath.Join(l.root, "00000 Джедайские техники/small.png")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://cover/large.png", filepath.Join(l.root, "00000 Джедайские техники/large.png")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://zip", filepath.Join(l.root, "00000 Джедайские техники/audiobook/zip/Джедайские техники.zip")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://pdf", filepath.Join(l.root, "00000 Джедайские техники/e-book/pdf/Джедайские техники.pdf")).Return(nil).Once()

	amk.On("DownloadFile", ctxtest.Match, "https://fb2", filepath.Join(l.root, "00000 Джедайские техники/e-book/fb2/Джедайские техники.fb2")).
		Return(&api.Error{Code: 404}).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://mobi", filepath.Join(l.root, "00000 Джедайские техники/e-book/mobi/Джедайские техники.mobi")).
		Return(&url.Error{Err: errors.New("stopped after 10 redirects")}).Once()

	ch := make(chan book.Book, 2)

	bk := book.Book{
		Title: "Джедайские техники",
		Photos: []book.Address{{
			URL: "https://photos/photos.png",
		}},
		Cover: book.Cover{
			Small: "https://cover/small.png",
			Large: "https://cover/large.png",
		},
		Files: book.Files{
			Books: map[string]book.Addresses{
				"pdf": {
					book.Address{
						URL: "https://pdf",
					},
				},
				"mobi": {
					book.Address{
						URL: "https://mobi",
					},
				},
				"fb2": {
					book.Address{
						URL: "https://fb2",
					},
				},
			},
			AudioBooks: map[string]book.Addresses{
				"zip": {
					book.Address{
						URL: "https://zip",
					},
				},
			},
			Demo: map[string]book.Addresses{
				"epub": {
					book.Address{
						URL: "https://epub",
					},
				},
			},
		},
	}
	ch <- bk
	// duplicate to check the protection from re-downloading.
	ch <- bk
	close(ch)

	require.NoError(t, l.Worker(ctxtest.Background(), ch))
	indexFile := filepath.Join(l.root, "00000 Джедайские техники/book.json")
	require.FileExists(t, indexFile)

	wantData, err := json.MarshalIndent(bk, "", "\t")
	require.NoError(t, err)

	fileData, err := ioutil.ReadFile(indexFile)
	require.NoError(t, err)
	require.JSONEq(t, string(wantData), string(fileData))
}

func Test_cutter(t *testing.T) {
	tests := map[string]string{}

	const maxFileNameLen = 255
	for i := maxFileNameLen; i < 1024; i++ {
		tests[genStaticFileName(i)] = genStaticFileName(maxFileNameLen)
	}

	for i := maxFileNameLen; i > 5; i-- {
		tests[genStaticFileName(i)] = genStaticFileName(i)
	}

	for arg, want := range tests {
		got := cutter(arg)
		require.Equalf(t, want, got, "want length: %d, got length: %d", len(want), len(got))
		require.Equal(t, filepath.Dir(want), filepath.Dir(got), "want dir: %q, got dir: %q")
	}

	t.Run("test on the example of a specific case for book 3344", func(t *testing.T) {
		golden.Equal(t, []byte(cutter(string(golden.Read(t)))))
	})
}

func genStaticFileName(size int) string {
	const prefix = "/var/folders/yh/x8t18v653t752p5nk400qt580000gn/T/"
	return prefix + strings.Repeat("w", size-4) + ".txt"
}

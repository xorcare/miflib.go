package downloader

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/xorcare/miflib.go/internal/api"
	"github.com/xorcare/miflib.go/internal/book"
	"github.com/xorcare/miflib.go/internal/book/files"
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

	amk.On("DownloadFile", ctx, "https://ogg/0", "jedi/audiobook/ogg/Введение.ogg").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://ogg/1", "jedi/audiobook/ogg/Глава 1.ogg").Return(nil).Once()
	amk.On("DownloadFile", ctx, "https://ogg/2", "jedi/audiobook/ogg/Приложения.ogg").Return(nil).Once()

	require.NoError(t, l.downloadAudiobook(ctx, "jedi", book.Book{
		Title: "Джедайские техники",
		Files: files.Files{
			AudioBooks: map[string]files.Addresses{
				"m4b": {
					files.Address{
						URL: "https://m4b",
					},
				},
				"mp3": {
					files.Address{
						URL:   "https://mp3/0",
						Title: "Введение",
					}, files.Address{
						URL:   "https://mp3/1",
						Title: "Глава 1",
					}, files.Address{
						URL:   "https://mp3/2",
						Title: "Приложения",
					},
				},
				"ogg": {
					files.Address{
						URL:   "https://ogg/0",
						Title: "Введение",
					}, files.Address{
						URL:   "https://ogg/1",
						Title: "Глава 1",
					}, files.Address{
						URL:   "https://ogg/2",
						Title: "Приложения",
					},
				},
				"zip": {
					files.Address{
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
			Files: files.Files{
				AudioBooks: map[string]files.Addresses{
					"zip": {
						files.Address{
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
		Files: files.Files{
			Books: map[string]files.Addresses{
				"epub": {
					files.Address{
						URL: "https://epub",
					},
				},
				"fb2": {
					files.Address{
						URL: "https://fb2",
					},
				},
				"mobi": {
					files.Address{
						URL: "https://mobi",
					},
				},
				"pdf": {
					files.Address{
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
			Files: files.Files{
				Books: map[string]files.Addresses{
					"epub": {
						files.Address{
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
		Files: files.Files{
			Demo: map[string]files.Addresses{
				"epub": {
					files.Address{
						URL: "https://epub",
					},
				},
				"fb2": {
					files.Address{
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
			Files: files.Files{
				Demo: map[string]files.Addresses{
					"epub": {
						files.Address{
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
		Photos: []files.Address{
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
			Photos: []files.Address{
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

	amk.On("DownloadFile", ctxtest.Match, "https://photos/photos.png", "jedi/photos/photos.png").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://cover/small.png", "jedi/small.png").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://cover/large.png", "jedi/large.png").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://zip", "jedi/audiobook/zip/Джедайские техники.zip").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://pdf", "jedi/e-book/pdf/Джедайские техники.pdf").Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://epub", "jedi/demo/epub/Джедайские техники.epub").Return(nil).Once()

	amk.On("DownloadFile", ctxtest.Match, "https://fb2", "jedi/e-book/fb2/Джедайские техники.fb2").
		Return(&api.Error{Code: 404}).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://mobi", "jedi/e-book/mobi/Джедайские техники.mobi").
		Return(&url.Error{Err: errors.New("stopped after 10 redirects")}).Once()

	require.NoError(t, l.download(ctxtest.Background(), "jedi", book.Book{
		Title: "Джедайские техники",
		Photos: []files.Address{{
			URL: "https://photos/photos.png",
		}},
		Cover: book.Cover{
			Small: "https://cover/small.png",
			Large: "https://cover/large.png",
		},
		Files: files.Files{
			Books: map[string]files.Addresses{
				"pdf": {
					files.Address{
						URL: "https://pdf",
					},
				},
				"mobi": {
					files.Address{
						URL: "https://mobi",
					},
				},
				"fb2": {
					files.Address{
						URL: "https://fb2",
					},
				},
			},
			AudioBooks: map[string]files.Addresses{
				"zip": {
					files.Address{
						URL: "https://zip",
					},
				},
			},
			Demo: map[string]files.Addresses{
				"epub": {
					files.Address{
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

	amk.On("DownloadFile", ctxtest.Match, "https://photos/photos.png", filepath.Join(l.root, "00000 Джедайские техники/photos/photos.png")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://cover/small.png", filepath.Join(l.root, "00000 Джедайские техники/small.png")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://cover/large.png", filepath.Join(l.root, "00000 Джедайские техники/large.png")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://zip", filepath.Join(l.root, "00000 Джедайские техники/audiobook/zip/Джедайские техники.zip")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://pdf", filepath.Join(l.root, "00000 Джедайские техники/e-book/pdf/Джедайские техники.pdf")).Return(nil).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://epub", filepath.Join(l.root, "00000 Джедайские техники/demo/epub/Джедайские техники.epub")).Return(nil).Once()

	amk.On("DownloadFile", ctxtest.Match, "https://fb2", filepath.Join(l.root, "00000 Джедайские техники/e-book/fb2/Джедайские техники.fb2")).
		Return(&api.Error{Code: 404}).Once()
	amk.On("DownloadFile", ctxtest.Match, "https://mobi", filepath.Join(l.root, "00000 Джедайские техники/e-book/mobi/Джедайские техники.mobi")).
		Return(&url.Error{Err: errors.New("stopped after 10 redirects")}).Once()

	ch := make(chan book.Book, 2)

	bk := book.Book{
		Title: "Джедайские техники",
		Photos: []files.Address{{
			URL: "https://photos/photos.png",
		}},
		Cover: book.Cover{
			Small: "https://cover/small.png",
			Large: "https://cover/large.png",
		},
		Files: files.Files{
			Books: map[string]files.Addresses{
				"pdf": {
					files.Address{
						URL: "https://pdf",
					},
				},
				"mobi": {
					files.Address{
						URL: "https://mobi",
					},
				},
				"fb2": {
					files.Address{
						URL: "https://fb2",
					},
				},
			},
			AudioBooks: map[string]files.Addresses{
				"zip": {
					files.Address{
						URL: "https://zip",
					},
				},
			},
			Demo: map[string]files.Addresses{
				"epub": {
					files.Address{
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

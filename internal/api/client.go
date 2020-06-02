package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/xorcare/miflib.go/internal/osutil"
)

var errClientDidNotAuthenticate = errors.New("client did not authenticate, please authenticate first")

var _ doer = &http.Client{}

type doer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type logger interface {
	Debugf(msg string, keysAndValues ...interface{})
}

// Client client for working with the miflib api.
type Client struct {
	http     doer
	basepath string

	log logger

	mx   sync.RWMutex
	once sync.Once

	authenticated bool
}

// Login is a authentication method.
func (c *Client) Login(ctx context.Context, username, password string) (err error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	loginURL := fmt.Sprintf("%s/auth/login.ajax", c.basepath)
	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBufferString(
		fmt.Sprintf(`{"email":%q,"password":%q}`, username, password),
	))
	if err != nil {
		return err
	}

	var res *http.Response
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if res, err = c.doRequest(ctx, req); err != nil {
		return err
	}
	defer res.Body.Close()

	c.authenticated = true

	return nil
}

// List is a method for getting a list of books.
func (c *Client) List(ctx context.Context) (resp ListResponse, err error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	if !c.authenticated {
		return resp, errClientDidNotAuthenticate
	}

	listURL := fmt.Sprintf("%s/books/list.ajax", c.basepath)
	req, err := http.NewRequest(http.MethodGet, listURL, nil)
	if err != nil {
		return ListResponse{}, err
	}

	var res *http.Response
	if res, err = c.doRequest(ctx, req); err != nil {
		return ListResponse{}, err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return ListResponse{}, err
	}

	return resp, nil
}

// DownloadFile this is the place to upload files.
func (c *Client) DownloadFile(ctx context.Context, url, filename string) (err error) {
	var req *http.Request
	var res *http.Response

	if exist, err := osutil.FileExists(filename); exist && err == nil {
		info, err := os.Stat(filename)
		if req, err = http.NewRequest(http.MethodHead, url, nil); err != nil {
			return err
		}
		if res, err = c.doRequest(ctx, req); err != nil {
			return err
		}
		defer res.Body.Close()
		if res.Header.Get("Content-Length") == strconv.FormatInt(info.Size(), 10) {
			c.log.Debugf("skip downloading url %q because file %q exist with equal size: %d",
				url, filename, info.Size())
			return nil
		}
	} else if err != nil {
		return err
	}

	c.mx.RLock()
	defer c.mx.RUnlock()
	if !c.authenticated {
		return errClientDidNotAuthenticate
	}
	filename, err = filepath.Abs(filename)
	if err != nil {
		return err
	}

	c.log.Debugf("downloading data from url %q to file %q", url, filename)

	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return err
	}

	if res, err = c.doRequest(ctx, req); err != nil {
		return err
	}
	defer res.Body.Close()

	if err = checkResponse(res); err != nil {
		return err
	}

	if err = os.MkdirAll(path.Dir(filename), 0755); err != nil {
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

func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	c.log.Debugf("http request is in progress, method: %q, url: %q", req.Method, req.URL)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(res); err != nil {
		defer res.Body.Close()
		return nil, err
	}

	return res, nil
}

// NewClient creates new instance of api client.
func NewClient(basepath string, logger logger, opts ...Option) *Client {
	c := &Client{
		basepath: basepath,
		http:     &http.Client{},
		log:      logger,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

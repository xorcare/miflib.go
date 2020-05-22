package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var errClientDidNotAuthenticate = errors.New("client did not authenticate, please authenticate first")

var _ doer = &http.Client{}

type doer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// Client client for working with the miflib api.
type Client struct {
	http     doer
	basepath string

	mx   sync.RWMutex
	once sync.Once

	authenticated bool
}

// Login is a authentication method.
func (c *Client) Login(ctx context.Context, username, password string) error {
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
	c.mx.RLock()
	defer c.mx.RUnlock()
	if !c.authenticated {
		return errClientDidNotAuthenticate
	}
	filename, err = filepath.Abs(filename)
	if err != nil {
		return err
	}

	log.Println("start download from url:", url, "to file:", filename)
	defer log.Println("finish download from url:", url, "to file:", filename)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := checkResponse(res); err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
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
func NewClient(basepath string, opts ...Option) *Client {
	c := &Client{
		basepath: basepath,
		http:     &http.Client{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

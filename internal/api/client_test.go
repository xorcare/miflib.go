package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xorcare/golden"
	"go.uber.org/zap"

	"github.com/xorcare/miflib.go/internal/ctxtest"
)

type doerMock struct {
	mock.Mock
}

func (d *doerMock) Do(req *http.Request) (resp *http.Response, err error) {
	a := d.Called(req)
	if a.Get(0).(*http.Response) == nil {
		return nil, a.Error(1)
	}
	return a.Get(0).(*http.Response), a.Error(1)
}

func TestClient_Login(t *testing.T) {
	dm := new(doerMock)
	defer dm.AssertExpectations(t)
	c := Client{
		http:     dm,
		basepath: "https://localhost:65535",
		log:      zap.NewNop().Sugar(),
	}

	dm.On("Do", mock.Anything).Run(checkRequest(t)).
		Return(httptest.NewRecorder().Result(), nil).Once()

	require.NoError(t, c.Login(ctxtest.Background(), t.Name(), t.Name()))
}

func TestClient_List(t *testing.T) {
	dm := new(doerMock)
	defer dm.AssertExpectations(t)
	c := Client{
		http:          dm,
		basepath:      "https://localhost:65535",
		log:           zap.NewNop().Sugar(),
		authenticated: true,
	}

	want := ListResponse{
		Total: 42,
	}
	resp := httptest.NewRecorder()
	require.NoError(t, json.NewEncoder(resp).Encode(want))

	dm.On("Do", mock.Anything).Run(checkRequest(t)).
		Return(resp.Result(), nil).Once()

	lr, err := c.List(ctxtest.Background())
	require.NoError(t, err)
	require.Equal(t, want, lr)
}

func TestClient_DownloadFile(t *testing.T) {
	dm := new(doerMock)
	defer dm.AssertExpectations(t)
	c := Client{
		http:          dm,
		log:           zap.NewNop().Sugar(),
		authenticated: true,
	}

	tempFile, err := ioutil.TempFile("", t.Name())
	require.NoError(t, err)
	tempFile.Close()
	require.NoError(t, os.Remove(tempFile.Name()))

	resp := httptest.NewRecorder()
	resp.Write([]byte(tempFile.Name()))
	dm.On("Do", mock.Anything).Run(checkRequest(t)).
		Return(resp.Result(), nil).Once()

	require.NoError(t, c.DownloadFile(ctxtest.Background(), "https://localhost:65535/tmp/favicon.json", tempFile.Name()))
	data, err := ioutil.ReadFile(tempFile.Name())
	require.NoError(t, err)
	require.Equal(t, tempFile.Name(), string(data))
}

func checkRequest(t *testing.T) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		req, ok := args[0].(*http.Request)
		require.Equal(t, true, ok, "first parameter should be *http.Request")
		require.Equal(t, true, ctxtest.Is(req.Context()), "should get ctxtest.Background()")
		dump, err := httputil.DumpRequest(req, true)
		require.NoError(t, err)
		golden.Equal(t, dump).FailNow()
	}
}

package httpoh

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	url string
	//body   io.Reader
	method string
}

func (r testRequest) Method() string { return r.method }
func (r testRequest) URL() string    { return r.url }

//func (r testRequest) Body() io.Reader { return r.body }

type testResponse struct {
	code int
	body []byte
}

func (resp *testResponse) ProcessResponse(r *http.Response) error {
	resp.code = r.StatusCode
	io.Copy(bytes.NewBuffer(resp.body), r.Body)
	return nil
}

func TestSuccessRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "httpoh", r.Header.Get("User-Agent"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := Config{
		UserAgent: "httpoh",
	}

	client, newError := NewClientNative(cfg, server.Client())
	require.NoError(t, newError)

	req := NewMockRequest(t)
	req.EXPECT().URL().Return(server.URL)
	req.EXPECT().Method().Return(http.MethodGet)

	resp := NewMockResponse(t)
	resp.EXPECT().ProcessResponse(mock.Anything).Run(func(r *http.Response) {
		assert.Equal(t, r.StatusCode, http.StatusOK)
		body := bytes.NewBuffer(make([]byte, 0, 100))
		io.Copy(body, r.Body)
		assert.Equal(t, body.String(), "OK")
	}).Return(nil)

	gotError := client.PerformRequest(context.Background(), req, resp)
	assert.NoError(t, gotError)
}

func TestPerformRequest(t *testing.T) {
	for _, tc := range []struct {
		Name              string
		Config            Config
		ServerHandler     http.Handler
		ServerStopped     bool
		RequestMethod     string
		RequestURL        string
		ResponseProcessor func(*http.Response) error
		WantErrorMatch    []string
	}{
		{
			Name: "get success",
			Config: Config{
				UserAgent: "custom ua",
			},
			RequestMethod: "GET",
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "custom ua", r.Header.Get("User-Agent"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			ResponseProcessor: func(r *http.Response) error {
				assert.Equal(t, r.StatusCode, http.StatusOK)
				body := bytes.NewBuffer(make([]byte, 0, 100))
				io.Copy(body, r.Body)
				assert.Equal(t, body.String(), "OK")
				return nil
			},
		},
		{
			Name:          "get success",
			Config:        Config{},
			RequestMethod: "GET",
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				expectDefaultUA := "Mozilla/5.0 (Linux; Android 11; Pixel 3a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.101 Mobile Safari/537.36"
				assert.Equal(t, expectDefaultUA, r.Header.Get("User-Agent"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			ResponseProcessor: func(r *http.Response) error {
				assert.Equal(t, r.StatusCode, http.StatusOK)
				body := bytes.NewBuffer(make([]byte, 0, 100))
				io.Copy(body, r.Body)
				assert.Equal(t, body.String(), "OK")
				return nil
			},
		},
		{
			Name:          "new request error",
			Config:        Config{},
			RequestMethod: "GET\n",
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			WantErrorMatch: []string{"invalid method"},
		},
		{
			Name:          "no server response",
			Config:        Config{},
			RequestMethod: "GET",
			ServerStopped: true,
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			WantErrorMatch: []string{"connection refused"},
		},
		{
			Name: "get response error",
			Config: Config{
				UserAgent: "custom ua",
			},
			RequestMethod: "GET",
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "custom ua", r.Header.Get("User-Agent"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			ResponseProcessor: func(r *http.Response) error {
				assert.Equal(t, r.StatusCode, http.StatusOK)
				body := bytes.NewBuffer(make([]byte, 0, 100))
				io.Copy(body, r.Body)
				assert.Equal(t, body.String(), "OK")
				return errors.New("WTF")
			},
			WantErrorMatch: []string{"WTF"},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("/", tc.ServerHandler)
			server := httptest.NewServer(mux)
			defer server.Close()

			client, newError := NewClientNative(tc.Config, server.Client())
			require.NoError(t, newError)

			req := NewMockRequest(t)
			req.EXPECT().URL().Return(server.URL + tc.RequestURL)
			req.EXPECT().Method().Return(tc.RequestMethod)

			resp := NewMockResponse(t)
			if tc.ResponseProcessor != nil {
				resp.EXPECT().ProcessResponse(mock.Anything).RunAndReturn(tc.ResponseProcessor)
			}

			if tc.ServerStopped {
				server.Close()
			}
			gotError := client.PerformRequest(context.Background(), req, resp)

			if tc.WantErrorMatch == nil {
				assert.NoError(t, gotError)
			} else if assert.Error(t, gotError) {
				gotErrorString := gotError.Error()
				for _, substr := range tc.WantErrorMatch {
					assert.Contains(t, gotErrorString, substr)
				}
			}
		})
	}
}

func TestRequestWithHeaders(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		Config         Config
		ServerHandler  http.Handler
		RequestMethod  string
		RequestHeaders http.Header
	}{
		{
			Name:          "two headers",
			RequestMethod: "POST",
			Config:        Config{UserAgent: "UA"},
			RequestHeaders: http.Header{
				"h1": []string{"v1", "v2"},
				"h2": []string{"v3"},
			},
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "UA", r.Header.Get("User-Agent"))
				assert.Equal(t, []string{"v1", "v2"}, r.Header.Values("h1"))
				assert.Equal(t, "v3", r.Header.Get("h2"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
		{
			Name:          "override useragent",
			RequestMethod: "POST",
			Config:        Config{UserAgent: "UA"},
			RequestHeaders: http.Header{
				"User-Agent": []string{"v1"},
				"h2":         []string{"v2"},
			},
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, []string{"v1"}, r.Header.Values("User-Agent"))
				assert.Equal(t, []string{"v2"}, r.Header.Values("h2"))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			server := httptest.NewServer(tc.ServerHandler)
			defer server.Close()

			req := NewMockRequestWithHeaders(t)
			req.EXPECT().URL().Return(server.URL)
			req.EXPECT().Method().Return(tc.RequestMethod)
			req.EXPECT().Headers().Return(tc.RequestHeaders)

			client, newError := NewClientNative(tc.Config, server.Client())
			require.NoError(t, newError)

			resp := NewMockResponse(t)
			resp.EXPECT().ProcessResponse(mock.Anything).RunAndReturn(func(r *http.Response) error {
				assert.Equal(t, r.StatusCode, http.StatusOK)
				body := bytes.NewBuffer(make([]byte, 0, 2))
				io.Copy(body, r.Body)
				assert.Equal(t, body.String(), "OK")
				return nil
			})

			gotError := client.PerformRequest(context.Background(), req, resp)
			assert.NoError(t, gotError)
		})
	}
}

func TestRequestWithBody(t *testing.T) {
	for _, tc := range []struct {
		Name          string
		Config        Config
		ServerHandler http.Handler
		RequestMethod string
		RequestBody   io.Reader
	}{
		{
			Name:          "post request with body",
			RequestMethod: "POST",
			Config:        Config{UserAgent: "UA"},
			RequestBody:   strings.NewReader("BODY"),
			ServerHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "UA", r.Header.Get("User-Agent"))

				body := bytes.NewBuffer(make([]byte, 0, 4))
				io.Copy(body, r.Body)
				assert.Equal(t, body.String(), "BODY")

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			server := httptest.NewServer(tc.ServerHandler)
			defer server.Close()

			req := NewMockRequestWithBody(t)
			req.EXPECT().URL().Return(server.URL)
			req.EXPECT().Method().Return(tc.RequestMethod)
			req.EXPECT().Body().Return(tc.RequestBody)

			client, newError := NewClientNative(tc.Config, server.Client())
			require.NoError(t, newError)

			resp := NewMockResponse(t)
			resp.EXPECT().ProcessResponse(mock.Anything).RunAndReturn(func(r *http.Response) error {
				assert.Equal(t, r.StatusCode, http.StatusOK)
				body := bytes.NewBuffer(make([]byte, 0, 2))
				io.Copy(body, r.Body)
				assert.Equal(t, body.String(), "OK")
				return nil
			})

			gotError := client.PerformRequest(context.Background(), req, resp)
			assert.NoError(t, gotError)
		})
	}
}

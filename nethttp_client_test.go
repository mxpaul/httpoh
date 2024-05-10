package httpoh

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := Config{}

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

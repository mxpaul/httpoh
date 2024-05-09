package httpoh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	req := testRequest{url: server.URL, method: http.MethodGet}
	resp := testResponse{}

	gotError := client.PerformRequest(context.Background(), req, &resp)
	fmt.Printf("After\n")
	assert.NoError(t, gotError)
	assert.Equal(t, resp.code, http.StatusOK)
	assert.Equal(t, string(resp.body), "OK")
}

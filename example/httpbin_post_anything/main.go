package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mxpaul/httpoh"
)

type HTTPBinResponse struct {
	Data    string            `json:"data,omitempty"`
	Origin  string            `json:"origin,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type AnythingRequest struct {
	url  string
	body string
}

func (req *AnythingRequest) URL() string    { return req.url }
func (req *AnythingRequest) Method() string { return http.MethodPost }
func (req *AnythingRequest) Headers() http.Header {
	return http.Header{"content-type": []string{"text/plain"}}
}
func (req *AnythingRequest) Body() io.Reader {
	return strings.NewReader(req.body)
}

type AnythingResponse struct {
	Code     int
	Response HTTPBinResponse
}

func (resp *AnythingResponse) ProcessResponse(r *http.Response) error {
	resp.Code = r.StatusCode

	err := json.NewDecoder(r.Body).Decode(&resp.Response)
	if err != nil {
		fmt.Errorf("%w: json parse error", err)
	}

	return nil
}

func main() {
	cfg := httpoh.Config{
		UserAgent:      "httpoh v0.1.1",
		ConnectTimeout: 1 * time.Second,
	}

	httpClient, err := httpoh.NewNetHTTPClient(cfg)
	if err != nil {
		panic(err)
	}

	client, err := httpoh.NewClientNative(cfg, httpClient)
	if err != nil {
		panic(err)
	}

	req := &AnythingRequest{url: "http://httpbin.org/anything", body: `request data`}
	resp := &AnythingResponse{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.PerformRequest(ctx, req, resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response data: %v\n", resp.Response.Data)
	fmt.Printf("Response headers: %v\n", resp.Response.Headers)
	fmt.Printf("Response origin: %v\n", resp.Response.Origin)
}

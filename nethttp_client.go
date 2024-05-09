package httpoh

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/Azure/go-ntlmssp"
)

func NewNetHTTPClient(cfg Config) (*http.Client, error) {
	dialer := net.Dialer{
		Timeout: cfg.ConnectTimeout,
	}

	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: cfg.TLSHandshakeTimeout,
		DisableCompression:  cfg.DisableCompression,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:     cfg.MaxConnsPerHost,
	}

	var clientTransport http.RoundTripper = transport
	if cfg.WithNTLM {
		clientTransport = ntlmssp.Negotiator{RoundTripper: transport}
	}

	c := &http.Client{
		Transport: clientTransport,
		Timeout:   cfg.ReadWriteTimeout,
	}

	if !cfg.FollowRedirect {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return c, nil
}

type ClientNative struct {
	HTTP      *http.Client
	UserAgent string
}

var _ Client = (*ClientNative)(nil)

func NewClientNative(cfg Config, httpClient *http.Client) (*ClientNative, error) {
	c := &ClientNative{
		HTTP:      httpClient,
		UserAgent: cfg.UserAgent,
	}
	return c, nil
}

func (c *ClientNative) PerformRequest(ctx context.Context, req Request, resp Response) error {

	netReq, _ := http.NewRequestWithContext(ctx, req.Method(), req.URL(), nil)
	netResp, _ := c.HTTP.Do(netReq)
	fmt.Printf("CODE: %d\n", netResp.StatusCode)

	_ = resp.ProcessResponse(netResp)

	return nil
}

package httpoh

import (
	"context"
	"crypto/tls"
	"io"
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
	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent()
	}
	return c, nil
}

func defaultUserAgent() string {
	return "Mozilla/5.0 (Linux; Android 11; Pixel 3a) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.101 Mobile Safari/537.36"
}

func (c *ClientNative) PerformRequest(ctx context.Context, req Request, resp Response) error {
	httpRequestMethod, httpRequestURL := req.Method(), req.URL()
	var httpRequestBody io.Reader
	if bReq, implements := req.(RequestWithBody); implements {
		httpRequestBody = bReq.Body()
	}

	netReq, err := http.NewRequestWithContext(ctx, httpRequestMethod, httpRequestURL, httpRequestBody)
	if err != nil {
		return err
	}

	netReq.Header.Set("User-Agent", c.UserAgent)
	if hReq, implements := req.(RequestWithHeaders); implements {
		for name, vals := range hReq.Headers() {
			for i, value := range vals {
				if _, exist := netReq.Header[name]; i == 0 && exist {
					netReq.Header.Set(name, value)
				} else {
					netReq.Header.Add(name, value)
				}
			}
		}
	}

	netResp, err := c.HTTP.Do(netReq)
	if err != nil {
		return err
	}
	defer netResp.Body.Close()

	err = resp.ProcessResponse(netResp)

	return err
}

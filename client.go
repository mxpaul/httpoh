package httpoh

import (
	"context"
	"io"
	"net/http"
)

type Client interface {
	PerformRequest(ctx context.Context, req Request, res Response) error
}

type Request interface {
	Method() string
	URL() string
}

type RequestWithHeaders interface {
	Request
	Headers() http.Header
}

type RequestWithBody interface {
	Request
	Body() io.Reader
}

type Response interface {
	ProcessResponse(r *http.Response) error
}

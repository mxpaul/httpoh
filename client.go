package httpoh

import (
	"context"
	"net/http"
)

type Client interface {
	PerformRequest(ctx context.Context, req Request, res Response) error
}

type Request interface {
	Method() string
	URL() string
}

//Body() io.Reader

type Response interface {
	ProcessResponse(r *http.Response) error
}

package httpoh

import "time"

type Config struct {
	UserAgent           string
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	ConnectTimeout      time.Duration
	ReadWriteTimeout    time.Duration
	TLSHandshakeTimeout time.Duration
	DisableCompression  bool
	FollowRedirect      bool
	InsecureSkipVerify  bool
	WithNTLM            bool
}

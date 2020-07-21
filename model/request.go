package model

import (
	"io"
	"strings"
	"time"
)

type Request interface {
	GetProtocol() string
}

type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string // Headers
	Body    string            // body
	Timeout time.Duration     // 请求超时时间
}

func (H *HTTPRequest) GetProtocol() string {
	return "http"
}

func (H *HTTPRequest) GetBody() io.Reader{
	return strings.NewReader(H.Body)
}

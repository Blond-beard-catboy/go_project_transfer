package proxy

import (
	"net/http/httputil"
	"net/url"
)

func NewProxy(targetURL string) *httputil.ReverseProxy {
	url, _ := url.Parse(targetURL)
	return httputil.NewSingleHostReverseProxy(url)
}

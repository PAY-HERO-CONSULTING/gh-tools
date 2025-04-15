package httpclient

import (
	"crypto/tls"
	"net/http"
)

type (
	HttpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

func NewHttpClient() HttpClient {

	return http.DefaultClient
}

func NewWithTLSConfig() HttpClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}

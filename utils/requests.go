package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/PAY-HERO-CONSULTING/gh-tools/httpclient"
)

type response struct {
	Body       []byte
	Header     http.Header
	StatusCode int
}

func NewJsonRequest(
	ctx context.Context,
	method string,
	endpoint string,
	requestBody any,
	headers map[string]string,
) (*http.Request, error) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return &http.Request{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return &http.Request{}, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}

func DoRequest(
	client httpclient.HttpClient,
	request *http.Request,
) (*response, error) {
	resp, err := client.Do(request)
	if err != nil {
		return &response{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response{}, err
	}

	response := &response{
		Body:       body,
		Header:     resp.Header,
		StatusCode: resp.StatusCode,
	}

	return response, nil
}

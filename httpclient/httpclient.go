package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func Request(url, method string, body map[string]interface{}) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("url is invalid")
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 299 || response.StatusCode < 200 {
		return nil, fmt.Errorf("request failed with status code %d, resp %s", response.StatusCode, string(bodyResp))
	}

	return bodyResp, nil
}

type HttpClient struct {
	URL     string
	Method  string
	Body    map[string]interface{}
	Header  map[string]string
	Timeout time.Duration
}

func New() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) SetURL(url string) *HttpClient {
	h.URL = url
	return h
}

func (h *HttpClient) SetMethod(method string) *HttpClient {
	h.Method = method
	return h
}

func (h *HttpClient) SetBody(body map[string]interface{}) *HttpClient {
	h.Body = body
	return h
}

func (h *HttpClient) SetTimeout(timeout time.Duration) *HttpClient {
	h.Timeout = timeout
	return h
}

func (h *HttpClient) SetHeader(headers map[string]string) *HttpClient {
	h.Header = headers
	return h
}

func (h *HttpClient) Execute() ([]byte, error) {
	if h.URL == "" {
		return nil, fmt.Errorf("url is invalid")
	}
	requestBody, err := json.Marshal(h.Body)
	if err != nil {
		return nil, err
	}
	if h.Timeout == 0 {
		h.Timeout = 120 * time.Second
	}
	client := &http.Client{
		Timeout: h.Timeout,
	}
	if h.Method == "" {
		h.Method = "GET"
	}
	request, err := http.NewRequest(h.Method, h.URL, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	if len(h.Header) > 0 {
		for k, v := range h.Header {
			request.Header.Set(k, v)
		}
	}
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bodyResp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 299 || response.StatusCode < 200 {
		return nil, fmt.Errorf("request failed with status code %d, resp %s", response.StatusCode, string(bodyResp))
	}

	return bodyResp, nil
}

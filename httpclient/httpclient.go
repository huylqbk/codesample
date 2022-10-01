package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
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

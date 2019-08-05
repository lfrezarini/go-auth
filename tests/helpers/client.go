package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/99designs/gqlgen/client"
)

type HTTPClient struct{}

// DoRequest is a helper that provides a http client with the possibility to pass a map of headers (The gqlgen client library doesn't provide a way to do it)
func (h HTTPClient) DoRequest(URL, query string, headers map[string]string) ([]byte, error) {
	httpClient := http.Client{}

	request := client.Request{
		Query: query,
	}

	requestBody, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to create request body: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/graphql", URL), bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, fmt.Errorf("Error while trying to create the request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to do the request: %v", err)
	}

	f, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to read the request body: %v", err)
	}

	resp.Body.Close()

	return f, nil
}

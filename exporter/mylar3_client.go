// Copyright 2024, Esteban Sanchez

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package exporter

import (
	"io"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type mylar3Client struct {
	client *http.Client
	opts   *Mylar3Opts
	logger *logrus.Logger
}

type Mylar3Opts struct {
	URI     string
	APIKey  string
	Timeout int
}
type Mylar3RawResponse struct {
	Data []byte
	Err  error
}

func newMylar3Client(opts *Mylar3Opts, logger *logrus.Logger) *mylar3Client {
	return &mylar3Client{
		client: &http.Client{},
		opts:   opts,
		logger: logger,
	}
}

func (c *mylar3Client) makeRequest(command string, params map[string]string) ([]byte, error) {
	request, err := http.NewRequest("GET", c.opts.URI, nil)

	if err != nil {
		c.logger.Errorf("Error getting %s: %s", c.opts.URI, err)
		return nil, err
	}
	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	query.Add("apikey", c.opts.APIKey)
	query.Add("cmd", command)
	request.URL.RawQuery = query.Encode()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		c.logger.Errorf("Error getting %s: %s", c.opts.URI, response.Status)
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("Error reading response body: %s", err)
		return nil, err
	}

	return body, nil
}

func (c *mylar3Client) CallCommand(command string, params map[string]string) ([]byte, error) {
	return c.makeRequest(command, params)
}

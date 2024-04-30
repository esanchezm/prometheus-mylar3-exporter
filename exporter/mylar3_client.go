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

func newMylar3Client(opts *Mylar3Opts, logger *logrus.Logger) *mylar3Client {
	return &mylar3Client{
		client: &http.Client{},
		opts:   opts,
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

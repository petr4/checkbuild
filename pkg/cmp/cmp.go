package cmp

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Client get http client
type Client struct {
	rst *resty.Client
}

type Result struct {
	Version     string `json:"version"`
	Build       string `json:"build"`
	Environment string `json:"environment"`
	StatusCode  string `json:"status"`
	Url         string `json:"url"`
}

// Init config
func Init() (*Client, error) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	if viper.GetBool("debug") {
		client.SetDebug(true)
	}
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "checkbuild-agent",
	})
	client.SetTimeout(10 * time.Second)

	return &Client{
		rst: client,
	}, nil

}

// Run compare urls
func (c *Client) Run(urls []string) ([]Result, bool, error) {
	var wg sync.WaitGroup
	rs := make(chan Result, 1)
	var results []Result
	for _, s := range urls {
		wg.Add(1)
		go func(s string, out chan Result) (bool, error) {
			defer wg.Done()

			var result Result
			resp, err := c.rst.R().ForceContentType("application/json").
				SetResult(&result).
				Get(s)
			result.StatusCode = fmt.Sprintf("%v", resp.StatusCode())
			result.Url = s
			logrus.Infof("Resp: %v", result)
			out <- result
			if err != nil {
				logrus.Errorf("Error executing GET: %v", err)
				return false, fmt.Errorf("Error executing GET on url: %v", "url")
			}
			return true, nil
		}(s, rs)

	}
	go func() {
		wg.Wait()
		close(rs)
	}()

	for r := range rs {
		results = append(results, r)
	}
	if len(results) <= 1 {
		logrus.Infof("Resp: %v", len(results))
		return nil, false, fmt.Errorf("To compare urls number must be >=2")
	}
	for i := 0; i < len(results)-1; i++ {
		if results[i].Build != results[i+1].Build || results[i].StatusCode != "200" {
			return results, false, nil
		}
	}

	return results, true, nil

}

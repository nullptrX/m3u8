package tool

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Get(url string) (io.ReadCloser, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	c := http.Client{
		Timeout: time.Duration(60) * time.Second,
		Transport: tr,
	}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}

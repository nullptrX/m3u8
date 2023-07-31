package tool

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/go-resty/resty/v2"
	"github.com/nullptrx/v2/common"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	urllib "net/url"
	"time"
)

func BuildClient() *http.Client {
	var transport *http.Transport
	if common.Proxy != "" {
		proxy := func(_ *http.Request) (*urllib.URL, error) {
			return urllib.Parse(common.Proxy)
		}
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: proxy,
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	// Cookie handle
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar:       jar,
		Timeout:   time.Second * 30,
		Transport: transport,
	}
}

func Get(c *http.Client, url string) (io.ReadCloser, error) {
	// Create a Resty Client
	client := resty.NewWithClient(c)
	//resp, err := client.R().
	//	SetHeaders(common.Headers).
	//	Head(url)
	resp, err := client.R().
		SetHeaders(common.Headers).
		Get(url)
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("http status code %d, error: %v", resp.StatusCode(), resp.Error())
	}
	var body io.Reader
	encoding := resp.Header().Get("Content-Encoding")
	reader := bytes.NewReader(resp.Body())
	if encoding == "gzip" {
		body, err = gzip.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("Failed to create gzip reader: %v", err)
		}
	} else if encoding == "br" {
		body = brotli.NewReader(reader)
	} else if encoding == "deflate" {
		body = flate.NewReader(reader)
	} else {
		body = ioutil.NopCloser(reader)
	}
	return io.NopCloser(body), nil
}

package tool

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/nullptrx/v2/common"
	"io"
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

func Get(client *http.Client, url string) (io.ReadCloser, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", common.UserAegnt)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8 coding\tgzip, deflate, br dnt\t1")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}
	var body io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		defer resp.Body.Close()
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to create gzip reader: %v", err)
		}
	} else {
		body = resp.Body
	}

	//content, _ := io.ReadAll(body)
	//_ = content
	return body, nil
}

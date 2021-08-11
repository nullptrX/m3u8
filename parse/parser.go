package parse

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/nullptrx/v2/tool"
	"io/ioutil"
	"net/url"
	"os"
)

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
	Keys map[int]string
}

func FromURL(link string, keyPath string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()
	body, err := tool.Get(link)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer body.Close()
	m3u8, err := parse(body)
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(tool.ResolveURL(u, sf.URI), keyPath)
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[int]string),
	}

	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == CryptMethodNONE:
			continue
		case key.Method == CryptMethodAES:
			// Request URL to extract decryption key
			var err error
			keyURL := key.URI
			keyURL = tool.ResolveURL(u, keyURL)
			var keyByte []byte
			keyByte, err = ioutil.ReadFile(keyPath)
			if err != nil {
				resp, err := tool.Get(keyURL)
				if err != nil {
					return nil, fmt.Errorf("extract key failed: %s", err.Error())
				}
				keyByte, err = ioutil.ReadAll(resp)
				_ = resp.Close()
			}
			fmt.Println("decryption key: ", base64.StdEncoding.EncodeToString(keyByte))
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}

func FromFile(filePath string, keyPath string) (*Result, error) {
	body, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer body.Close()
	m3u8, err := parse(body)
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(sf.URI, keyPath)
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		M3u8: m3u8,
		Keys: make(map[int]string),
	}

	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == CryptMethodNONE:
			continue
		case key.Method == CryptMethodAES:
			// Request URL to extract decryption key
			var err error
			keyURL := key.URI
			var keyByte []byte
			keyByte, err = ioutil.ReadFile(keyPath)
			if err != nil {
				resp, err := tool.Get(keyURL)
				if err != nil {
					return nil, fmt.Errorf("extract key failed: %s", err.Error())
				}
				keyByte, err = ioutil.ReadAll(resp)
				_ = resp.Close()
			}
			fmt.Println("decryption key: ", base64.StdEncoding.EncodeToString(keyByte))
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}

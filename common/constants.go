package common

var (

	// socks://127.0.0.1:6153
	Proxy   = ""
	Headers = map[string]string{
		"User-Aegnt":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate, br",
	}
)

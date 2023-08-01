package main

import (
	"flag"
	"fmt"
	"github.com/nullptrx/v2/common"
	"github.com/nullptrx/v2/dl"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	nurl "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	url      string
	output   string
	chanSize int
	verbose  bool
	key      string
	merge    bool
	proxy    string
	config   string
)

func init() {
	flag.StringVar(&url, "u", "", "URL, required")
	flag.IntVar(&chanSize, "n", 10, "Maximum number of occurrences")
	flag.StringVar(&output, "o", strconv.FormatInt(time.Now().Unix(), 10), "Output folder, required")
	flag.BoolVar(&verbose, "v", false, "Verbose log, optional")
	flag.StringVar(&key, "k", "", "Key path, optional")
	flag.BoolVar(&merge, "m", false, "Merge files, optional")
	flag.StringVar(&proxy, "p", "", "Proxy url (such as socks://127.0.0.1:1080, http://127.0.0.1:1080), optional")
	flag.StringVar(&config, "c", "dump.yaml", "Config file for http headers.")
}

func main() {
	flag.Parse()
	u, err := nurl.Parse(url)
	if err != nil {
		panicParameter("u")
	}

	file, err := ioutil.ReadFile(config)
	if err == nil {
		var config map[string]string
		err = yaml.Unmarshal(file, &config)
		if err == nil {
			if config["Referer"] == "" {
				config["Referer"] = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
			}
			for key, value := range config {
				common.Headers[key] = value
			}
			if verbose {
				for key, value := range common.Headers {
					fmt.Printf("%s: %s\n", key, value)
				}
			}
		} else {
			if verbose {
				fmt.Println("[warning]", err)
			}
		}
	} else {
		if verbose {
			fmt.Println("[warning]", err)
		}
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	common.Proxy = proxy
	if !strings.HasPrefix(url, "http") {
		if len(flag.Args()) > 0 {
			for _, arg := range flag.Args() {
				if strings.HasPrefix(arg, "http") {
					url = arg
					break
				}
			}
		}
		if url == "" {
			panicParameter("u")
		}
	}
	//if output == "" {
	//	panicParameter("o")
	//}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}

	isM3u8 := strings.Contains(u.Path, ".m3u8")
	if isM3u8 {
		downloader, err := dl.NewTask(output, url, verbose, key)
		if err != nil {
			panic(err)
		}
		if merge {
			if err := downloader.Merge(); err != nil {
				panic(err)
			}
		} else {
			if err := downloader.Start(chanSize); err != nil {
				panic(err)
			}
		}
		fmt.Println("Done!")
	} else {
		dl.DirectDownload(output, url, chanSize, verbose)
	}

}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}

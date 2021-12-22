package main

import (
	"flag"
	"fmt"
	"github.com/nullptrx/v2/dl"
	nurl "net/url"
	"os"
	"strings"
)

var (
	url      string
	output   string
	chanSize int
	verbose  bool
	key      string
	merge    bool
)

func init() {
	flag.StringVar(&url, "u", "", "URL, required")
	flag.IntVar(&chanSize, "c", 10, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
	flag.BoolVar(&verbose, "v", false, "Verbose log, optional")
	flag.StringVar(&key, "k", "", "Key path, optional")
	flag.BoolVar(&merge, "m", false, "Merge files, optional")
}

func main() {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if url == "" {
		if len(flag.Args()) > 0 {
			url = flag.Arg(0)
		} else {
			panicParameter("u")
		}
	}
	//if output == "" {
	//	panicParameter("o")
	//}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}

	u, err := nurl.Parse(url)
	if err != nil {
		panicParameter("u")
	}

	isM3u8 := strings.HasSuffix(u.Path, ".m3u8")

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

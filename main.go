package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nullptrx/v2/dl"
)

var (
	url      string
	output   string
	chanSize int
	verbose  bool
	key      string
)

func init() {
	flag.StringVar(&url, "u", "", "M3U8 URL, required")
	flag.IntVar(&chanSize, "c", 10, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
	flag.BoolVar(&verbose, "v", false, "Verbose log, optional")
	flag.StringVar(&key, "k", "", "Key path, optional")
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
		panicParameter("u")
	}
	if output == "" {
		//panicParameter("o")
		output = strconv.FormatInt(time.Now().Unix(), 10)
	}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}
	downloader, err := dl.NewTask(output, url, verbose, key)
	if err != nil {
		panic(err)
	}
	if err := downloader.Start(chanSize); err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}

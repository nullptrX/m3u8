package dl

import (
	"fmt"
	"github.com/nullptrx/v2/common"
	"github.com/nullptrx/v2/link/pkg/base"
	"github.com/nullptrx/v2/link/pkg/download"
	//"github.com/monkeyWie/gopeed-core/pkg/base"
	//"github.com/monkeyWie/gopeed-core/pkg/download"
	"github.com/nullptrx/v2/m3u8/tool"

	nurl "net/url"
	"os"
	"path/filepath"
)

func DirectDownload(output, url string, chansize int, verbose bool) {

	//var wg sync.WaitGroup
	//
	// wg.Done()
	//
	// wg.Add(1)
	// wg.Wait()

	var folder string
	// If no output folder specified, use current directory
	if output == "" {
		current, err := tool.CurrentDir()
		if err != nil {
			fmt.Printf("Download fail:%v\n", err)
			return
		}
		folder = filepath.Join(current, output)
	} else {
		folder = output
	}
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		fmt.Printf("Download fail:%v\n", err)
		return
	}
	u, _ := nurl.Parse(url)
	filename := filepath.Base(u.Path)

	finallyCh := make(chan error)
	err := download.Boot().
		URL(url).
		Extra(&base.Extra{
			Header: map[string]string{
				"User-Agent": common.UserAegnt,
			},
		}).
		Listener(func(event *download.Event) {
			if event.Key == download.EventKeyProgress {
				printProgress("Downloading", event)
			}
			if event.Key == download.EventKeyFinally {
				printProgress(" Completed ", event)
				finallyCh <- event.Err
			}

		}).
		Create(&base.Options{
			Name:        filename,
			Path:        folder,
			Connections: chansize,
		})
	if err != nil {
		panic(err)
	}
	err = <-finallyCh
	if err != nil {
		fmt.Printf("\nDownload fail:%v\n", err)
	} else {
		fmt.Println("\nDone!")
	}
}

func printProgress(title string, event *download.Event) {
	task := event.Task
	rate := float32(task.Progress.Downloaded) / float32(task.Res.TotalSize)
	tool.DrawProgressBar(title, rate, progressWidth)
}

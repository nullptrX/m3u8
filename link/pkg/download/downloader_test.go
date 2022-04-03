package download

import (
	"github.com/nullptrx/v2/link/internal/protocol/http"
	"github.com/nullptrx/v2/link/pkg/base"
	test2 "github.com/nullptrx/v2/link/pkg/test"
	"reflect"
	"sync"
	"testing"
)

func TestDownloader_Resolve(t *testing.T) {
	listener := test2.StartTestFileServer()
	defer listener.Close()

	downloader := NewDownloader(http.FetcherBuilder)
	req := &base.Request{
		URL: "http://" + listener.Addr().String() + "/" + test2.BuildName,
	}
	res, err := downloader.Resolve(req)
	if err != nil {
		t.Fatal(err)
	}
	want := &base.Resource{
		Req:       req,
		TotalSize: test2.BuildSize,
		Range:     true,
		Files: []*base.FileInfo{
			{
				Name: test2.BuildName,
				Path: "",
				Size: test2.BuildSize,
			},
		},
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Resolve error = %s, want %s", test2.ToJson(res), test2.ToJson(want))
	}
}

func TestDownloader_Create(t *testing.T) {
	listener := test2.StartTestFileServer()
	defer listener.Close()

	downloader := NewDownloader(http.FetcherBuilder)
	req := &base.Request{
		URL: "http://" + listener.Addr().String() + "/" + test2.BuildName,
	}
	res, err := downloader.Resolve(req)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	downloader.Listener(func(event *Event) {
		if event.Key == EventKeyDone {
			wg.Done()
		}
	})
	err = downloader.Create(res, &base.Options{
		Path:        test2.Dir,
		Name:        test2.DownloadName,
		Connections: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	want := test2.FileMd5(test2.BuildFile)
	got := test2.FileMd5(test2.DownloadFile)
	if want != got {
		t.Errorf("Download error = %v, want %v", got, want)
	}
}

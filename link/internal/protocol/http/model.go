package http

import (
	"github.com/nullptrx/v2/link/pkg/base"
)

type Chunk struct {
	Status     base.Status
	Begin      int64
	End        int64
	Downloaded int64
}

func NewChunk(begin int64, end int64) *Chunk {
	return &Chunk{
		Status: base.DownloadStatusReady,
		Begin:  begin,
		End:    end,
	}
}

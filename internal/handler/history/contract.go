package history

import (
	"context"
	"time"
)

type requester interface {
	RequestHistory(context.Context, int, time.Time) (string, error)
}

type downloader interface {
	DownloadFile(context.Context, string) ([]byte, error)
}

type historian interface {
	requester
	downloader
}

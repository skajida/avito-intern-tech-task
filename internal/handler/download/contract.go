package download

import "context"

type downloader interface {
	DownloadFile(context.Context, string) ([]byte, error)
}

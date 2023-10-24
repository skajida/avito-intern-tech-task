//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package download

import (
	"context"
	mod "service-segs/internal/model"
)

type downloader interface {
	DownloadFile(context.Context, mod.Filename) (mod.RawData, error)
}

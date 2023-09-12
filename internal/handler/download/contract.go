//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package download

import "context"

type downloader interface {
	DownloadFile(context.Context, string) ([]byte, error)
}

//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package history

import (
	"context"
	mod "service-segs/internal/model"
	"time"
)

type requester interface {
	RequestHistory(context.Context, int, time.Time) (mod.Filename, error)
}

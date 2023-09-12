//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package history

import (
	"context"
	"time"
)

type requester interface {
	RequestHistory(context.Context, int, time.Time) (string, error)
}

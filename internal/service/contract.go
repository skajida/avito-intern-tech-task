//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package service

import (
	"context"
	"service-segs/internal/model/exchange"
	"time"
)

type erepository interface {
	Exists(context.Context, int) bool
	Count(context.Context) int
}

type tsegments interface {
	InsertSegment(context.Context, string) error
	DeleteSegment(context.Context, string) error
}

type tbelongings interface {
	SelectBelonging(context.Context, int) ([]string, error)
	UpdateBelonging(context.Context, int, []string, []string) error
	UpdateBelongingTimer(context.Context, int, []string, []string, time.Time) error
	SelectHistory(context.Context, int, time.Time, time.Time) ([]exchange.HistoryEntry, error)
}

type irepository interface {
	tsegments
	tbelongings
}

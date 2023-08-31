//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package service

import (
	"context"
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

type HistoryEntry struct {
	UserId    int       `csv:"user_id"`
	SegTag    string    `csv:"seg_id"`
	Operation string    `csv:"operation"`
	Time      time.Time `csv:"timestamp"`
}

type tbelongings interface {
	SelectBelonging(context.Context, int) ([]string, error)
	UpdateBelonging(context.Context, int, []string, []string) error
	SelectHistory(context.Context, int, time.Time, time.Time) ([]HistoryEntry, error)
}

type irepository interface {
	tsegments
	tbelongings
}

//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package service

import (
	"context"
	mod "service-segs/internal/model"
	"time"
)

type csvRepository interface {
	CreateHistoryFile(context.Context, mod.HistoryCollection) (mod.Filename, error)
	DownloadHistoryFile(context.Context, mod.Filename) (mod.RawData, error)
}

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
	SelectHistory(context.Context, int, time.Time, time.Time) (mod.HistoryCollection, error)
	AutoApply(context.Context, string, []int) error
}

type irepository interface {
	tsegments
	tbelongings
}

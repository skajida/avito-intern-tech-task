//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package service

import "context"

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
}

type irepository interface {
	tsegments
	tbelongings
}

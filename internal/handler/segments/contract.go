//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package segments

import "context"

type creator interface {
	AddSegment(context.Context, string) error
	AddSegmentAutoApply(context.Context, string, float64) error
}

type deletor interface {
	DeleteSegment(context.Context, string) error
}

type segmentor interface {
	creator
	deletor
}

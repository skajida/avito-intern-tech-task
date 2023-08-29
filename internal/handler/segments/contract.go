//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package segments

import (
	"context"
	"service-segs/internal/model/exchange/request"
)

type creator interface {
	AddSegment(context.Context, request.Segment) error
}

type deletor interface {
	RemoveSegment(context.Context, request.Segment) error
}

type segmentor interface {
	creator
	deletor
}

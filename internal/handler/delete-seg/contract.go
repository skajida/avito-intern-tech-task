//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package delete_seg

import "context"

type service interface {
	RemoveSegment(ctx context.Context) error
}

//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package create_seg

import "context"

type service interface {
	AddSegment(ctx context.Context) error
}

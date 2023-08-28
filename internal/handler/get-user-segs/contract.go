//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package get_user_segs

import "context"

type service interface {
	GetUserSegments(ctx context.Context) error
}

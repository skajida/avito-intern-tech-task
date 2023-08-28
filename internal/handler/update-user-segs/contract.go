//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package update_user_segs

import "context"

type service interface {
	UpdateUserSegments(ctx context.Context) error
}

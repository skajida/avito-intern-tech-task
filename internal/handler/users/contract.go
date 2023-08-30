//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package users

import "context"

type getter interface {
	GetUserSegments(context.Context, int) ([]string, error)
}

type updater interface {
	UpdateUserSegments(context.Context, int, []string, []string) error
}

type useror interface {
	getter
	updater
}

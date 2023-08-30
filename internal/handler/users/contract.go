//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package users

import (
	"context"
)

type getter interface {
	GetUserSegments(context.Context, uint) ([]string, error)
}

type updater interface {
	UpdateUserSegments(context.Context, uint, []string, []string) error
}

type useror interface {
	getter
	updater
}

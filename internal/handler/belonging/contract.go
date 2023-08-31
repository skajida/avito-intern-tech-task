//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package belonging

import "context"

type reader interface {
	ReadBelonging(context.Context, int) ([]string, error)
}

type modifier interface {
	ModifyBelonging(context.Context, int, []string, []string) error
}

type belonger interface {
	reader
	modifier
}

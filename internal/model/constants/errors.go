package constants

import "errors"

var (
	WrongRequest   = errors.New("wrong request format")
	WrongUser      = errors.New("wrong user format")
	NotFound       = errors.New("not found")
	InvalidSegment = errors.New("invalid segment")
	InternalError  = errors.New("internal error")
)

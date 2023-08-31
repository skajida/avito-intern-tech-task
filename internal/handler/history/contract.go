package history

import (
	"context"
	"time"
)

type requester interface {
	RequestHistory(context.Context, int, time.Time) (string, error)
}

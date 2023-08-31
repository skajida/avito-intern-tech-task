package exchange

import "time"

type HistoryEntry struct {
	UserId    int       `csv:"user_id"`
	SegTag    string    `csv:"seg_id"`
	Operation string    `csv:"operation"`
	Time      time.Time `csv:"timestamp"`
}

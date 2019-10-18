package ledger

import "time"

type RecordFilter func(*Record) *Record

type DateRangeFilter struct {
	Since time.Time
	Until time.Time
}

func (f *DateRangeFilter) Filter(r *Record) *Record {
	if r.Date.Before(f.Since) || r.Date.After(f.Until) {
		return nil
	}

	return r
}

package datetime

import "time"

type Month struct {
	year  int
	month int
}

func NewMonth(t time.Time) Month {
	return Month{year: t.Year(), month: int(t.Month())}
}

func (m Month) Same(other Month) bool {
	return m.year == other.year && m.month == other.month
}

func (m Month) Before(other Month) bool {
	if m.year < other.year {
		return true
	}

	if m.year == other.year && m.month <= other.month {
		return true
	}

	return false
}

func (m Month) After(other Month) bool {
	return !m.Before(other)
}

func (m Month) String() string {
	startOfMonth := time.Date(m.year, time.Month(m.month), 1, 0, 0, 0, 0, time.UTC)
	return startOfMonth.Format("Jan 2006")
}

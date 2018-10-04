package datetime

import (
	"fmt"
	"time"
)

const dateLayout = "02 Jan 06"

// Week presents a week.
type Week struct {
	year int
	week int
}

func (w Week) Same(other Week) bool {
	return w.year == other.year && w.week == other.week
}

func (w Week) Before(other Week) bool {
	if w.year < other.year {
		return true
	}

	if w.year == other.year && w.week < other.week {
		return true
	}

	return false
}

func (w Week) After(other Week) bool {
	return !w.Before(other)
}

func (w Week) Dates() (time.Time, time.Time) {
	return CommercialDate(w.year, w.week, 1),
		CommercialDate(w.year, w.week, 7)
}

func (w Week) String() string {
	startOfWeek, endOfWeek := w.Dates()
	return fmt.Sprintf("%s - %s", startOfWeek.Format(dateLayout), endOfWeek.Format(dateLayout))
}

func NewWeek(t time.Time) Week {
	year, week := t.ISOWeek()
	return Week{year: year, week: week}
}

package datetime

import (
	"time"
)

const daysInWeek = 7

// BeginningOfWeek returns a new date/time which represents the start of this week.
func BeginningOfWeek(t time.Time) time.Time {
	currentDayNumber := 6
	if t.Weekday() != 0 {
		currentDayNumber = int(t.Weekday()) - 1
	}
	daysToWeekStart := currentDayNumber % 7

	diff := time.Hour*24*time.Duration(-daysToWeekStart) +
		time.Hour*time.Duration(-t.Hour()) +
		time.Minute*time.Duration(-t.Minute()) +
		time.Second*time.Duration(-t.Second()) +
		time.Nanosecond*time.Duration(-t.Nanosecond())

	return t.Add(diff)
}

// CommercialDate returns a date for a given year, week and weekday. The month and
// weekday values may be outside their ranges and will be normalized during the
// conversion. For example weekday 8 becomes 7.
func CommercialDate(year, week, weekday int) time.Time {
	base := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	baseYear, _ := base.ISOWeek()
	if baseYear == year {
		week = week - 1
	}

	baseWeekday := int(base.Weekday())
	if baseWeekday == 0 {
		baseWeekday = 7
	}

	// TODO: implement week and weekday normalization.
	offset := daysInWeek*week + (weekday - baseWeekday)

	return base.AddDate(0, 0, offset)
}

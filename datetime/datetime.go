package datetime

import "time"

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

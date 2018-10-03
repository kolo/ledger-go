package datetime

import (
	"testing"
	"time"
)

func Test_BeginningOfWeek(t *testing.T) {
	expected := _time("2018-05-28T00:00:00+02:00")

	tests := []string{
		"2018-05-28T21:51:17.113+02:00", // Monday
		"2018-05-29T21:51:17.113+02:00", // Tuesday
		"2018-05-30T21:51:17.113+02:00", // Wednesday
		"2018-05-31T21:51:17.113+02:00", // Thursday
		"2018-06-01T21:51:17.113+02:00", // Friday
		"2018-06-02T21:51:17.113+02:00", // Saturday
		"2018-06-03T21:51:17.113+02:00", // Sunday
	}

	for _, tc := range tests {
		assertEqualTime(t, expected, BeginningOfWeek(_time(tc)))
	}
}

func Test_CommercialDate(t *testing.T) {
	base := _time("2017-01-15T00:00:00Z") // Sunday
	for wd := 1; wd <= 7; wd++ {
		assertEqualTime(t, base.AddDate(0, 0, wd), CommercialDate(2017, 3, wd))
	}
}

func assertEqualTime(t *testing.T, expected, actual time.Time) {
	if expected.Equal(actual) {
		return
	}
	t.Errorf("\nError: time assertion failed\n\texpected: %v\n\t  actual: %v\n", expected, actual)
}

func _time(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

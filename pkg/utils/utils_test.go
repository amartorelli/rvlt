package utils

import (
	"testing"
	"time"
)

func TestDaysUntilBirthday(t *testing.T) {
	tt := []struct {
		bd  string
		now string
		res int
		err error
		msg string
	}{
		{"2006-04-02", "2016-04-01", 1, nil, "one day difference should return 1"},
		{"2006-04-05", "2016-04-01", 4, nil, "four day difference should return 4"},
		{"2116-04-05", "2016-04-01", 0, errDateComp, "a birthday in the future should give an error"},
		{"1996-02-29", "2019-02-27", 1, nil, "when born on a leap year we should celebrate on the 28th"},
	}

	for _, tc := range tt {
		tbd, _ := time.Parse("2006-01-02", tc.bd)
		tnow, _ := time.Parse("2006-01-02", tc.now)

		days, err := DaysUntilBirthday(tbd, tnow)
		if err != tc.err {
			t.Errorf("%s, returned error %+v, expected %+v", tc.msg, err, tc.err)
		}
		if days != tc.res {
			t.Errorf("%s, returned %d", tc.msg, days)
		}
	}
}

func TestIsLeapYear(t *testing.T) {
	tt := []struct {
		year   int
		isLeap bool
		msg    string
	}{
		{1996, true, "1996 should be a leap year"},
		{1994, false, "1994 should not be a leap year"},
		{2020, true, "2020 should be a leap year"},
		{2001, false, "2001 should not be a leap year"},
		{2004, true, "2004 should be a leap year"},
	}

	for _, tc := range tt {
		if IsLeapYear(tc.year) != tc.isLeap {
			t.Error(tc.msg)
		}
	}
}

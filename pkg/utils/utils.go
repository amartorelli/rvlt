package utils

import (
	"errors"
	"math"
	"time"
)

var (
	errDateComp = errors.New("error comparing dates")
)

// IsLeapYear returns true if the year is a leap year
func IsLeapYear(y int) bool {
	year := time.Date(y, time.December, 31, 0, 0, 0, 0, time.Local)
	days := year.YearDay()

	if days > 365 {
		return true
	}

	return false
}

// DaysUntilBirthday calculates the number of days until next birthday
func DaysUntilBirthday(bd, now time.Time) (int, error) {
	if bd.After(now) {
		return 0, errDateComp
	}

	var leapAdjustment int
	// if the person was born on a leap year we celebrate on the 28th of Feb
	if bd.Day() == 29 && bd.Month() == 02 && (!IsLeapYear(now.Year())) {
		leapAdjustment = -1
	}

	birthday := bd.AddDate(now.Year()-bd.Year(), 0, leapAdjustment)

	return int(math.Round(birthday.Sub(now).Hours() / 24)), nil
}
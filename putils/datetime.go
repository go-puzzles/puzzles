// File:		datetime.go
// Created by:	Hoven
// Created on:	2024-09-19
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package putils

import "time"

// BeginOfMinute return beginning minute time of day.
func BeginOfMinute(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, t.Location())
}

// EndOfMinute return end minute time of day.
func EndOfMinute(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfHour return beginning hour time of day.
func BeginOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
}

// EndOfHour return end hour time of day.
func EndOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfDay return beginning hour time of day.
func BeginOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay return end time of day.
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfWeek return beginning week, default week begin from Sunday.
func BeginOfWeek(t time.Time, beginFrom ...time.Weekday) time.Time {
	var beginFromWeekday = time.Sunday
	if len(beginFrom) > 0 {
		beginFromWeekday = beginFrom[0]
	}
	y, m, d := t.AddDate(0, 0, int(beginFromWeekday-t.Weekday())).Date()
	beginOfWeek := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	if beginOfWeek.After(t) {
		return beginOfWeek.AddDate(0, 0, -7)
	}
	return beginOfWeek
}

// EndOfWeek return end week time, default week end with Saturday.
func EndOfWeek(t time.Time, endWith ...time.Weekday) time.Time {
	var endWithWeekday = time.Saturday
	if len(endWith) > 0 {
		endWithWeekday = endWith[0]
	}
	y, m, d := t.AddDate(0, 0, int(endWithWeekday-t.Weekday())).Date()
	var endWithWeek = time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
	if endWithWeek.Before(t) {
		endWithWeek = endWithWeek.AddDate(0, 0, 7)
	}
	return endWithWeek
}

// BeginOfMonth return beginning of month.
func BeginOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth return end of month.
func EndOfMonth(t time.Time) time.Time {
	return BeginOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// BeginOfYear return the date time at the begin of year.
func BeginOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear return the date time at the end of year.
func EndOfYear(t time.Time) time.Time {
	return BeginOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// IsLeapYear check if param year is leap year or not.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// BetweenSeconds returns the number of seconds between two times.
func BetweenSeconds(t1 time.Time, t2 time.Time) int64 {
	index := t2.Unix() - t1.Unix()
	return index
}

// DayOfYear returns which day of the year the parameter date `t` is.
func DayOfYear(t time.Time) int {
	y, m, d := t.Date()
	firstDay := time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
	nowDate := time.Date(y, m, d, 0, 0, 0, 0, t.Location())

	return int(nowDate.Sub(firstDay).Hours() / 24)
}

// IsWeekend checks if passed time is weekend or not.
func IsWeekend(t time.Time) bool {
	return time.Saturday == t.Weekday() || time.Sunday == t.Weekday()
}

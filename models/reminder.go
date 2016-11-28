package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kathrynsiegel/reminder-bot/helpers"
)

// RepeatInterval type
type RepeatInterval int64

// Valid repeat interval consts
const (
	RepeatDaily RepeatInterval = iota
	RepeatWeekly
	RepeatMonthly
)

// Reminder represents one reminder a user has asked
// the bot to remember.
type Reminder struct {
	gorm.Model
	UserID            uint
	Recurring         bool
	RepeatInterval    RepeatInterval
	RepeatDay         time.Weekday // Only used for weekly repeats
	RepeatDayOfMonth  int          // Only used for monthly repeats
	RepeatTimeOfDayMs int64        // Time of day
	RepeatEvery       int64        // Repeat every x days, weeks, months
	Description       string
	NextSendAtMs      int64
	LastSendAtMs      int64
	Timezone          string
}

const durationDay = time.Hour * 24

// CalculateNextSendAtMs calculates the next time that a reminder
// should be sent. If it should never be sent again, NextSendAtMs
// is set to -1.
func (r *Reminder) CalculateNextSendAtMs() int64 {
	if r.Recurring == false && r.LastSendAtMs != 0 {
		// Reminder was sent and is not recurring.
		return -1
	}
	lastSendAtMs := r.NextSendAtMs
	if lastSendAtMs > 0 {
		switch r.RepeatInterval {
		case RepeatDaily:
			return lastSendAtMs + (r.RepeatEvery * helpers.DurationDayMs)
		case RepeatWeekly:
			return lastSendAtMs + (r.RepeatEvery * helpers.DurationWeekMs)
		case RepeatMonthly:
			return dateInNextMonth(lastSendAtMs, r.RepeatEvery)
		}
	}
	timeNow := time.Now()
	dayStart := timeNow.Truncate(durationDay)
	nextSendAt := dayStart.Add(time.Duration(r.RepeatTimeOfDayMs * int64(time.Millisecond)))
	if nextSendAt.Before(timeNow) {
		nextSendAt.AddDate(0, 0, 1)
	}
	switch r.RepeatInterval {
	case RepeatWeekly:
		// Increment to correct day of week
		for nextSendAt.Weekday() != r.RepeatDay {
			nextSendAt.Add(durationDay)
		}
	case RepeatMonthly:
		for nextSendAt.Day() != r.RepeatDayOfMonth {
			nextSendAt.Add(durationDay)
		}
	}
	return helpers.TimeToMs(nextSendAt)
}

func dateInNextMonth(prevSendMs int64, numSkippedMonths int64) int64 {
	prevSendTime := helpers.MsToTime(prevSendMs)
	return helpers.TimeToMs(prevSendTime.AddDate(0, int(numSkippedMonths), 0))
}

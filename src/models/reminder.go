package models

import "github.com/jinzhu/gorm"

// Day type
type Day int64

// Valid day constants
const (
	Sunday Day = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
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
	UserID            uint64
	Recurring         bool
	RepeatInterval    RepeatInterval
	RepeatDay         Day   // Only used for weekly repeats
	RepeatDayOfMonth  int64 // Only used for monthly repeats
	RepeatTimeOfDayMs int64 // Time of day
	RepeatEvery       int64 // Repeat every x days, weeks, months
}

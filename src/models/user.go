package models

type User struct {
	Id         int64
	FacebookId string
}

type Day int64

const (
	Sunday Day = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type RepeatInterval int64

const (
	RepeatDaily RepeatInterval = iota
	RepeatWeekly
	RepeatMonthly
)

type Reminder struct {
	Id                int64
	UserId            int64
	Recurring         bool
	RepeatInterval    RepeatInterval
	RepeatDay         Day   // Only used for weekly repeats
	RepeatDayOfMonth  int64 // Only used for monthly repeats
	RepeatTimeOfDayMs int64 // Time of day
	RepeatEvery       int64 // Repeat every x days, weeks, months
}

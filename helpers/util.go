package helpers

import "time"

// PanicIfError panics if passed an error.
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// TimeToMs converts a time object to its equivalent timestamp in ms.
func TimeToMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// MsToTime converts a timestamp in ms to a time object.
func MsToTime(timeMs int64) time.Time {
	return time.Unix(0, timeMs*int64(time.Millisecond))
}

// TimeNowInMs returns the current time in milliseconds.
func TimeNowInMs() int64 {
	return TimeToMs(time.Now())
}

const DurationDayMs = int64(24 * 60 * 60 * 1000)
const DurationWeekMs = DurationDayMs * 7

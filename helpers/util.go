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

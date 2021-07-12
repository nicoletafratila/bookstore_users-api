package date_utils

import (
	"time"
)

const (
	apiDbLayout = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowDbFormat() string {
	return GetNow().Format(apiDbLayout)
}

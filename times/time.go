package times

import (
	"time"
)

// 获取某一天的N点整点时间
func GetHourTime(d time.Time, hour int) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), hour, 0, 0, 0, d.Location())
}

// 获取当前时间的小时
func GetCurrentTimeHour(loc *time.Location) int {
	return time.Now().In(loc).Hour()
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetHourTime(d, 0)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

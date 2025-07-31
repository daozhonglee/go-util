package times

import (
	"time"
)

// DateTime64StringFromTimestampMilli 将毫秒时间戳转换为带毫秒的日期时间字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02 15:04:05.000" 的日期时间字符串
func DateTime64StringFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now + int64(tz_seconds)*1000
	tm := time.Unix(0, localtime*int64(time.Millisecond))
	return tm.In(time.UTC).Format(DateTimeMilliFormat)
}

// DateStringFromTimestampMilli 将毫秒时间戳转换为日期字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02" 的日期字符串
func DateStringFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now/1000 + int64(tz_seconds)
	return time.Unix(localtime, 0).In(time.UTC).Format(DateFormat)
}

// YesterdayStringFromTimestampMilli 将毫秒时间戳转换为昨天的日期字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02" 的昨天日期字符串
func YesterdayStringFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now/1000 + int64(tz_seconds)
	return time.Unix(localtime, 0).In(time.UTC).AddDate(0, 0, -1).Format(DateFormat)
}

// DateStringYMDHMSFromTimestampMilli 将毫秒时间戳转换为完整的日期时间字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02 15:04:05" 的日期时间字符串
func DateStringYMDHMSFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now/1000 + int64(tz_seconds)
	return time.Unix(localtime, 0).In(time.UTC).Format(DateTimeFormat)
}

// CurrentHourStringFromTimestampMilli 将毫秒时间戳转换为当前小时的开始时间字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02 15:04:05" 的当前小时开始时间字符串
func CurrentHourStringFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now/1000 + int64(tz_seconds)
	tm := time.Unix(localtime, 0).In(time.UTC)

	hourStart := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), 0, 0, 0, time.Local)
	dt := hourStart.Format(DateTimeFormat)
	return dt
}

// LastMondayStringFromTimestampMilli 将毫秒时间戳转换为上周一的日期字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02" 的上周一日期字符串
func LastMondayStringFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now/1000 + int64(tz_seconds)
	tm := time.Unix(localtime, 0).In(time.UTC)

	offset := (tm.Weekday() + 6) % 7

	l1time := tm.AddDate(0, 0, -int(offset)-7)

	dt := l1time.Format(DateFormat)
	return dt
}

// MondayStringFromTimestampMilli 将毫秒时间戳转换为本周一的日期字符串
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 格式为 "2006-01-02" 的本周一日期字符串
func MondayStringFromTimestampMilli(now int64, tz_seconds int) string {
	localtime := now/1000 + int64(tz_seconds)
	tm := time.Unix(localtime, 0).In(time.UTC)

	offset := int(time.Monday - tm.Weekday())
	if offset > 0 {
		offset = -6
	}

	l1time := tm.AddDate(0, 0, offset)

	dt := l1time.Format(DateFormat)
	return dt
}

// DateStringToTimestampMilli 将日期时间字符串转换为毫秒时间戳
// 参数:
//   - dt: 格式为 "2006-01-02 15:04:05" 的日期时间字符串
//
// 返回: 毫秒级时间戳，解析失败返回0
func DateStringToTimestampMilli(dt string) int64 {
	tm, err := time.ParseInLocation(DateTimeFormat, dt, time.UTC)
	if err != nil {
		return 0
	}
	return tm.Unix() * 1000
}

// TimestampMilliToTime 将毫秒时间戳转换为 time.Time 对象
// 参数:
//   - t: 毫秒级时间戳
//
// 返回: time.Time 对象
func TimestampMilliToTime(t int64) time.Time {
	return time.Unix(t/1000, 0)
}

// DateInt64FromTimestamp 将时间戳转换为当前小时开始的毫秒时间戳
// 流水小时统计使用【存储对应小时的毫秒值】
// 参数:
//   - now: 毫秒级时间戳
//
// 返回: 当前小时开始时间的毫秒时间戳
func DateInt64FromTimestamp(now int64) int64 {
	in := time.Unix(now/1000, 0).In(time.UTC)
	return time.Date(in.Year(), in.Month(), in.Day(), in.Hour(), 0, 0, 0, in.Location()).Unix() * 1000
}

// DateInt64FromTimestampMilli 将毫秒时间戳转换为当前小时开始的毫秒时间戳（带时区偏移）
// 参数:
//   - now: 毫秒级时间戳
//   - tz_seconds: 时区偏移量（秒）
//
// 返回: 当前小时开始时间的毫秒时间戳
func DateInt64FromTimestampMilli(now int64, tz_seconds int) int64 {
	localtime := now + int64(tz_seconds)*1000
	in := time.Unix(0, localtime*int64(time.Millisecond)).In(time.UTC)
	return time.Date(in.Year(), in.Month(), in.Day(), in.Hour(), 0, 0, 0, in.Location()).Unix() * 1000
}

// GetCurrentUnix 获取当前的Unix时间戳（秒级）
// 返回: 当前时间的秒级时间戳
func GetCurrentUnix() int64 {
	return time.Now().Unix()
}

// GetCurrentMilliUnix 获取当前的毫秒级时间戳
// 返回: 当前时间的毫秒级时间戳
func GetCurrentMilliUnix() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetCurrentNanoUnix 获取当前的纳秒级时间戳
// 返回: 当前时间的纳秒级时间戳
func GetCurrentNanoUnix() int64 {
	return time.Now().UnixNano()
}

// GetTwoDateDays 计算两个日期之间相差的天数
// 参数:
//   - startTime: 开始日期字符串，格式为 "2006-01-02"
//   - endTime: 结束日期字符串，格式为 "2006-01-02"
//   - loc: 时区位置
//
// 返回: 相差的天数（结束日期 - 开始日期）
func GetTwoDateDays(startTime, endTime string, loc *time.Location) int64 {
	startUnix, _ := time.ParseInLocation(DateFormat, startTime, loc)
	endUnix, _ := time.ParseInLocation(DateFormat, endTime, loc)
	// 求相差天数
	date := (endUnix.Unix() - startUnix.Unix()) / 86400
	return date
}

// Package util AT-GoLib 工具库主入口
// 重新导出各子包的常用函数，提供便捷的统一访问方式
package util

import (
	"github.com/daozhonglee/go-util/api"
	"github.com/daozhonglee/go-util/async"
	"github.com/daozhonglee/go-util/collection"
	"github.com/daozhonglee/go-util/crypto"
	"github.com/daozhonglee/go-util/delaytask"
	"github.com/daozhonglee/go-util/errorutil"
	"github.com/daozhonglee/go-util/id"
	"github.com/daozhonglee/go-util/json"
	"github.com/daozhonglee/go-util/random"
	times "github.com/daozhonglee/go-util/times"
	"github.com/daozhonglee/go-util/xss"
)

// 重新导出常用函数，用户可以选择直接使用 util.xxx 或 import 具体子包

// ID相关
var (
	ValidateChineseID = id.ValidateChinese
	CheckAge16To18    = id.CheckAge16To18
)

// 加密相关
var (
	MD5Hash            = crypto.MD5Hash
	SHA256Hash         = crypto.SHA256Hash
	HMACSign           = crypto.HMACSign
	NewMD5             = crypto.NewMD5
	InvisibleEncrypter = crypto.DefaultInvisibleEncrypter
)

// JSON相关
var (
	JSONMarshal   = json.Marshal
	JSONUnmarshal = json.Unmarshal
	JSONPretty    = json.Pretty
	JSONClean     = json.Clean
)

// 随机数相关
var (
	RandomInt = random.Int
)

// 时间相关
var (
	RelativeTime         = times.Relative
	GetCurrentUnix       = times.GetCurrentUnix
	GetCurrentMilliUnix  = times.GetCurrentMilliUnix
	GetCurrentNanoUnix   = times.GetCurrentNanoUnix
	TimestampMilliToTime = times.TimestampMilliToTime
	GetTwoDateDays       = times.GetTwoDateDays
)

// 时间格式常量
const (
	DateTimeFormat      = times.DateTimeFormat      // "2006-01-02 15:04:05"
	DateTimeMilliFormat = times.DateTimeMilliFormat // "2006-01-02 15:04:05.000"
	DateFormat          = times.DateFormat          // "2006-01-02"
)

// 异步相关
var (
	AsyncGo      = async.Go
	AsyncTimeout = async.Timeout
	AsyncSafe    = async.Safe
)

// API响应相关
var (
	APISuccess     = api.Success
	APIError       = api.Error
	APISuccessData = api.SuccessWithData
	APISuccessPage = api.SuccessWithPage
)

// 集合相关
var (
	NewSet = collection.NewSet
)

// XSS相关
var (
	XSSClean = xss.Clean
)

// 错误处理相关
var (
	PanicIf    = errorutil.PanicIf
	PanicStack = errorutil.PanicWithStack
	Recover    = errorutil.Recover
)

// 延时任务相关
var (
	DelayTaskPush         = delaytask.PushTask
	DelayTaskPull         = delaytask.PullTask
	DelayTaskPushInternal = delaytask.PushTaskInternal
	DelayTaskPullInternal = delaytask.PullTaskInternal
)

// 延时任务常量
const (
	DelayTaskIntervalMilliseconds = delaytask.INTERVAL_MILLISECONDS
	DelayTaskIntervalSeconds      = delaytask.INTERVAL_SECONDS
	DelayTaskIntervalMinutes      = delaytask.INTERVAL_MUNITES
	DelayTaskIntervalHour         = delaytask.INTERVAL_HOUR
)

package asyncio

import (
	"time"
)

// 重试函数
//
// times: 重试次数 负数则无限制
//
// delay: 休眠秒数 每次重试间休眠时间
//
// f: 要重试的函数
func Retry(times int, delay float64, f func() bool) bool {
	for ; times != 0 && !f(); times-- {
		time.Sleep(time.Duration(1000*delay) * time.Millisecond)
	}
	return times != 0
}

// 重试函数 支持参数
func RetryTask(times int, delay float64, task *Task) bool {
	return Retry(times, delay, func() bool { return task.Run().Bool() })
}

// 重试函数 通过是否抛出 error 判断
func RetryError(times int, delay float64, f func() error) bool {
	return Retry(times, delay, func() bool { return f() == nil })
}

// 重试函数 通过是否抛出 error 判断 支持参数
func RetryErrorTask(times int, delay float64, task *Task) bool {
	return Retry(times, delay, func() bool { return task.Run().Error() == nil })
}

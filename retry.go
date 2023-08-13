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
func Retry(times, delay int, f func() bool) {
	for ; times != 0 && !f(); times-- {
		if times > 0 {
			println("剩余重试次数:", times-1)
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// 重试函数 支持参数
func RetryWith[T any](times, delay int, coro Coro) {
	task := coro.ToTask()
	Retry(times, delay, func() bool { return task.first().(bool) })
}

// 重试函数 通过是否抛出 error 判断
func RetryError(times, delay int, f func() error) {
	Retry(times, delay, func() bool { return f() == nil })
}

// 重试函数 通过是否抛出 error 判断 支持参数
func RetryErrorWith(times, delay int, coro Coro) {
	task := coro.ToTask()
	Retry(times, delay, func() bool { return task.first().(error) == nil })
}

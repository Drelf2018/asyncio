package asyncio

import "sync"

var global EventLoop = NewEventLoop()

var NewEventLoop = func() EventLoop {
	return &AbstractEventLoop{
		status: Stop,
		tasks:  sync.WaitGroup{},
	}
}

func GetEventLoop() EventLoop {
	if global == nil {
		global = NewEventLoop()
	}
	return global
}

func SetEventLoop(fn func() EventLoop) {
	NewEventLoop = fn
	global = fn()
}

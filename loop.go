package asyncio

import (
	"sync"
)

var global *AbstractEventLoop

func NewEventLoop() *AbstractEventLoop {
	loop := new(AbstractEventLoop)
	loop.canRun.Add(1)
	return loop
}

func SetEventLoop(loop *AbstractEventLoop) {
	// close the old one
	if global != nil {
		global.tasks.Done()
	}
	// start the new one
	global = loop
	go global.RunForever()
}

func GetEventLoop() *AbstractEventLoop {
	if global == nil {
		panic("There has no eventloop.")
	}
	return global
}

type AbstractEventLoop struct {
	tasks  sync.WaitGroup
	canRun sync.WaitGroup
}

func (loop *AbstractEventLoop) run(task *Task) {
	loop.canRun.Wait()
	task.run()
	loop.tasks.Done()
}

func (loop *AbstractEventLoop) Coro(coro Coro) *Handle {
	task := coro.ToTask()
	loop.tasks.Add(1)
	go loop.run(task)
	return task.handle
}

func (loop *AbstractEventLoop) CreateTask(f any, args ...any) *Handle {
	return loop.Coro(Coro{f, args})
}

func (loop *AbstractEventLoop) RunUntilComplete() {
	loop.canRun.Done()
	loop.tasks.Wait()
}

func (loop *AbstractEventLoop) RunForever() {
	loop.tasks.Add(1)
	loop.RunUntilComplete()
}

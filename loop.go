package asyncio

import (
	"sync"
)

var global *AbstractEventLoop

func NewEventLoop() *AbstractEventLoop {
	loop := new(AbstractEventLoop)
	loop.frozen.Add(1)
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
		panic("There is no current eventloop.")
	}
	return global
}

type AbstractEventLoop struct {
	tasks  sync.WaitGroup
	frozen sync.WaitGroup
}

func (loop *AbstractEventLoop) Run(task *Task) {
	defer loop.tasks.Done()
	loop.frozen.Wait()
	task.Run()
}

func (loop *AbstractEventLoop) Coro(f any, args ...any) *Handle {
	return loop.CreateTask(Coro{f, args})
}

func (loop *AbstractEventLoop) CreateTask(coro Coro) *Handle {
	task := coro.Task()
	loop.tasks.Add(1)
	go loop.Run(task)
	return task.handle
}

func (loop *AbstractEventLoop) RunUntilComplete() {
	loop.frozen.Done()
	loop.tasks.Wait()
}

func (loop *AbstractEventLoop) RunForever() {
	loop.tasks.Add(1)
	loop.RunUntilComplete()
}

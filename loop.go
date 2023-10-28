package asyncio

import (
	"reflect"
	"sync"
)

type EventLoop interface {
	CreateTask(fn any, args ...any) *Task
	Run(task *Task)

	Start()
	RunUntilComplete(tasks ...*Task)

	RunForever(tasks ...*Task)
	Stop()
}

type Status int

const (
	Stop Status = iota
	Running
)

type AbstractEventLoop struct {
	status Status
	tasks  sync.WaitGroup
}

func (loop *AbstractEventLoop) CreateTask(fn any, args ...any) (r *Task) {
	task := &Task{Func: reflect.ValueOf(fn)}

	if len(args) == 2 && args[0] == reflect.Slice {
		arg := reflect.ValueOf(args[1])
		if arg.Kind() == reflect.Slice {
			r = task.Reflect(arg)
		}
	}

	if r == nil {
		r = task.Default(args...)
	}

	if loop.status == Running {
		loop.Run(r)
	}
	return
}

func (loop *AbstractEventLoop) Run(task *Task) {
	loop.tasks.Add(1)
	go func() {
		defer loop.tasks.Done()
		task.Run()
	}()
}

func (loop *AbstractEventLoop) Start() {
	loop.status = Running
}

func (loop *AbstractEventLoop) Finish() {
	loop.status = Stop
}

func (loop *AbstractEventLoop) RunUntilComplete(tasks ...*Task) {
	loop.Start()
	defer loop.Finish()
	for _, task := range tasks {
		loop.Run(task)
	}
	loop.tasks.Wait()
}

func (loop *AbstractEventLoop) RunForever(tasks ...*Task) {
	loop.tasks.Add(1)
	loop.RunUntilComplete(tasks...)
}

func (loop *AbstractEventLoop) Stop() {
	loop.tasks.Done()
}

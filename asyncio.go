package asyncio

import (
	"reflect"
	"slices"
)

type Args []any

func SingleArg[T any](args ...T) []Args {
	r := make([]Args, len(args))
	for i, arg := range args {
		r[i] = Args{arg}
	}
	return r
}

type Coro struct {
	Func any
	Args Args
}

func (c Coro) ToTask() *Task {
	r := make([]reflect.Value, len(c.Args))
	for i, arg := range c.Args {
		r[i] = reflect.ValueOf(arg)
	}
	return &Task{reflect.ValueOf(c.Func), r, new(Handle)}
}

func C(f any, args ...any) Coro {
	return Coro{f, args}
}

func NoArgsFunc(fs ...any) []Coro {
	r := make([]Coro, len(fs))
	for i, f := range fs {
		r[i] = Coro{Func: f}
	}
	return r
}

func Await(coros ...Coro) H {
	loop := NewEventLoop()
	r := make(H, len(coros))
	for i, coro := range coros {
		r[i] = loop.CreateTask(coro)
	}
	loop.RunUntilComplete()
	return r
}

func Slice(f any, args []Args) H {
	loop := NewEventLoop()
	r := make(H, len(args))
	for i, arg := range args {
		r[i] = loop.Coro(f, arg...)
	}
	loop.RunUntilComplete()
	return r
}

func List(f any, args []Args) H {
	loop := NewEventLoop()
	r := make(H, len(args))
	for i, arg := range args {
		r[i] = loop.Coro(f, slices.Insert(arg, 0, any(i))...)
	}
	loop.RunUntilComplete()
	return r
}

func Map[M ~map[K]V, K comparable, V any](f any, m M) H {
	loop := NewEventLoop()
	r := make(H, len(m))
	i := 0
	for k, v := range m {
		r[i] = loop.Coro(f, k, v)
		i++
	}
	loop.RunUntilComplete()
	return r
}

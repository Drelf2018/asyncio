package asyncio

import (
	"reflect"

	"golang.org/x/exp/slices"
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

func (c Coro) Task() *Task {
	return c.Callback(nil)
}

func (c Coro) Callback(f func([]any)) *Task {
	r := make([]reflect.Value, len(c.Args))
	for i, arg := range c.Args {
		r[i] = reflect.ValueOf(arg)
	}
	return &Task{reflect.ValueOf(c.Func), r, new(Handle), f}
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

func Do(f func(loop *AbstractEventLoop)) {
	loop := NewEventLoop()
	f(loop)
	loop.RunUntilComplete()
}

func Wait(coros ...Coro) H {
	r := make(H, len(coros))
	Do(func(loop *AbstractEventLoop) {
		for i, l := 0, len(coros); i < l; i++ {
			r[i] = loop.CreateTask(coros[i])
		}
	})
	return r
}

func ForEach[T any](args []T, f func(T)) {
	Do(func(loop *AbstractEventLoop) {
		for i, l := 0, len(args); i < l; i++ {
			loop.Coro(f, args[i])
		}
	})
}

func ForEachP[T any](args []T, f func(*T)) {
	Do(func(loop *AbstractEventLoop) {
		for i, l := 0, len(args); i < l; i++ {
			loop.Coro(f, &args[i])
		}
	})
}

func Slice(args []Args, f any) H {
	r := make(H, len(args))
	Do(func(loop *AbstractEventLoop) {
		for i, l := 0, len(args); i < l; i++ {
			r[i] = loop.Coro(f, args[i]...)
		}
	})
	return r
}

func List(args []Args, f any) H {
	r := make(H, len(args))
	Do(func(loop *AbstractEventLoop) {
		for i, arg := range args {
			r[i] = loop.Coro(f, slices.Insert(arg, 0, any(i))...)
		}
	})
	return r
}

func Map[M ~map[K]V, K comparable, V any](m M, f any) H {
	r := make(H, 0, len(m))
	Do(func(loop *AbstractEventLoop) {
		for k, v := range m {
			r = append(r, loop.Coro(f, k, v))
		}
	})
	return r
}

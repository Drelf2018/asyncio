package asyncio

import "reflect"

type Args []any

func SingleArg(args ...any) []Args {
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
	return &Task{reflect.ValueOf(c.Func), r}
}

func NoArgsFunc(fs ...any) []Coro {
	r := make([]Coro, len(fs))
	for i, f := range fs {
		r[i] = Coro{Func: f}
	}
	return r
}

func Await(coros ...Coro) []*Handle {
	loop := NewEventLoop()
	r := make([]*Handle, len(coros))
	for i, coro := range coros {
		r[i] = loop.CreateTask(coro)
	}
	loop.RunUntilComplete()
	return r
}

func Slice(f any, args ...Args) []*Handle {
	loop := NewEventLoop()
	r := make([]*Handle, len(args))
	for i, arg := range args {
		r[i] = loop.CreateTask(f, arg...)
	}
	loop.RunUntilComplete()
	return r
}

func Map[M ~map[K]V, K comparable, V any](f any, m M) []*Handle {
	loop := NewEventLoop()
	r := make([]*Handle, len(m))
	i := 0
	for k, v := range m {
		r[i] = loop.CreateTask(f, Args{k, v})
		i++
	}
	loop.RunUntilComplete()
	return r
}

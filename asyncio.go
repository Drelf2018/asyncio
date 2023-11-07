package asyncio

import (
	"sync"
)

func CreateTask(fn any, args ...any) *Task {
	return GetEventLoop().CreateTask(fn, args...)
}

func RunUntilComplete(tasks ...*Task) {
	GetEventLoop().RunUntilComplete(tasks...)
}

func RunForever(tasks ...*Task) {
	GetEventLoop().RunForever(tasks...)
}

func WaitGroup(delta int, f func(done func())) {
	wg := &sync.WaitGroup{}
	wg.Add(delta)
	f(wg.Done)
	wg.Wait()
}

func Loop(f func(loop EventLoop)) {
	loop := NewEventLoop()
	loop.Start()
	f(loop)
	loop.RunUntilComplete()
}

func Wait(tasks ...*Task) {
	NewEventLoop().RunUntilComplete(tasks...)
}

func Delay(seconds float64, fn any, args ...any) {
	go CreateTask(fn, args...).Delay(seconds)
}

func ForFunc[E any](arg E, f ...func(E)) {
	l := len(f)
	WaitGroup(l, func(done func()) {
		exec := func(i int) {
			defer done()
			f[i](arg)
		}
		for i := 0; i < l; i++ {
			go exec(i)
		}
	})
}

func ForFuncV[E any, V any](arg E, f ...func(E) V) []V {
	l := len(f)
	v := make([]V, l)
	WaitGroup(l, func(done func()) {
		exec := func(i int) {
			defer done()
			v[i] = f[i](arg)
		}
		for i := 0; i < l; i++ {
			go exec(i)
		}
	})
	return v
}

func ForEach[S ~[]E, E any](args S, f func(E)) {
	l := len(args)
	WaitGroup(l, func(done func()) {
		exec := func(i int) {
			defer done()
			f(args[i])
		}
		for i := 0; i < l; i++ {
			go exec(i)
		}
	})
}

func ForEachV[S ~[]E, E any, V any](args S, f func(E) V) []V {
	l := len(args)
	v := make([]V, l)
	WaitGroup(l, func(done func()) {
		exec := func(i int) {
			defer done()
			v[i] = f(args[i])
		}
		for i := 0; i < l; i++ {
			go exec(i)
		}
	})
	return v
}

func ForEachPtr[S ~[]E, E any](args S, f func(*E)) {
	l := len(args)
	WaitGroup(l, func(done func()) {
		exec := func(i int) {
			defer done()
			f(&args[i])
		}
		for i := 0; i < l; i++ {
			go exec(i)
		}
	})
}

func ForEachPtrV[S ~[]E, E any, V any](args S, f func(*E) V) []V {
	l := len(args)
	v := make([]V, l)
	WaitGroup(l, func(done func()) {
		exec := func(i int) {
			defer done()
			v[i] = f(&args[i])
		}
		for i := 0; i < l; i++ {
			go exec(i)
		}
	})
	return v
}

func Map[M ~map[K]V, K comparable, V any](m M, f func(K, V)) {
	WaitGroup(len(m), func(done func()) {
		exec := func(k K, v V) {
			defer done()
			f(k, v)
		}
		for k, v := range m {
			go exec(k, v)
		}
	})
}

func MapR[M ~map[K]V, K comparable, V any, R any](m M, f func(K, V) R) []R {
	l := len(m)
	r := make([]R, 0, l)
	WaitGroup(l, func(done func()) {
		exec := func(k K, v V) {
			defer done()
			r = append(r, f(k, v))
		}
		for k, v := range m {
			go exec(k, v)
		}
	})
	return r
}

func Slice[S ~[]E, E any](args S, f any) (tasks []*Task) {
	task := CreateTask(f)
	for _, arg := range args {
		if in, ok := any(arg).([]any); ok {
			tasks = append(tasks, task.Copy(in...))
		} else {
			tasks = append(tasks, task.Copy(arg))
		}
	}
	Wait(tasks...)
	return
}

func Fill[S ~[]E, E any](tasks []*Task, s S, position ...int) {
	p := 0
	if len(position) != 0 {
		p = position[0]
	}
	for idx, task := range tasks {
		s[idx] = task.result[p].(E)
	}
}

func Results[E any](tasks []*Task, position ...int) []E {
	s := make([]E, len(tasks))
	Fill(tasks, s, position...)
	return s
}

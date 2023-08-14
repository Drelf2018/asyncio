package asyncio

import (
	"reflect"
)

type Task struct {
	Func   reflect.Value
	Args   []reflect.Value
	handle *Handle
}

func (t *Task) run() {
	for _, v := range t.Func.Call(t.Args) {
		t.handle.out = append(t.handle.out, v.Interface())
	}
	t.handle.done = true
	t.handle.length = len(t.handle.out)
}

func (t Task) first() any {
	return t.Func.Call(t.Args)[0].Interface()
}

type Handle struct {
	done   bool
	length int
	out    []any
}

func (h *Handle) Done() bool {
	return h.done
}

func (h *Handle) Len() int {
	return h.length
}

func (h *Handle) Result() []any {
	return h.out
}

type H[T any] []*Handle

func (hs H[T]) To() []T {
	l := 0
	for _, h := range hs {
		l += h.Len()
	}
	rs := make([]T, l)
	i := 0
	for _, h := range hs {
		for _, r := range h.Result() {
			rs[i] = r.(T)
			i++
		}
	}
	return rs
}

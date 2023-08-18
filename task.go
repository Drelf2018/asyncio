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
}

func (t Task) first() any {
	return t.Func.Call(t.Args)[0].Interface()
}

type Handle struct {
	done bool
	out  []any
}

func (h *Handle) Done() bool {
	return h.done
}

func (h *Handle) Len() int {
	return len(h.out)
}

func (h *Handle) Result() []any {
	return h.out
}

type H []*Handle

func (H H) Len() (sum int) {
	for _, h := range H {
		sum += h.Len()
	}
	return
}

func To[S ~[]T, T any](H H, ts ...S) []T {
	var t []T
	if len(ts) == 0 {
		t = make([]T, H.Len())
	} else {
		t = ts[0]
	}
	i := 0
	for _, h := range H {
		for _, r := range h.Result() {
			t[i] = r.(T)
			i++
		}
	}
	return t
}

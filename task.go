package asyncio

import (
	"reflect"
)

type Task struct {
	Func     reflect.Value
	Args     []reflect.Value
	handle   *Handle
	callback func([]any)
}

func (t *Task) Run() {
	for _, v := range t.Func.Call(t.Args) {
		t.handle.out = append(t.handle.out, v.Interface())
	}
	t.handle.done = true
	if t.callback != nil {
		t.callback(t.handle.out)
	}
}

func (t *Task) Bool() bool {
	t.Run()
	if len(t.handle.out) == 0 {
		return false
	}
	return t.handle.out[0].(bool)
}

func (t *Task) Error() error {
	t.Run()
	if len(t.handle.out) == 0 {
		return nil
	}
	return t.handle.out[0].(error)
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

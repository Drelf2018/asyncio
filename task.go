package asyncio

import (
	"reflect"
)

type Task struct {
	Func reflect.Value
	Args []reflect.Value
}

func (t Task) run() (out []any) {
	for _, v := range t.Func.Call(t.Args) {
		out = append(out, v.Interface())
	}
	return
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

func (h *Handle) Result() []any {
	return h.out
}

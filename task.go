package asyncio

import (
	"errors"
	"reflect"
	"time"
)

var (
	ErrNoResult       = errors.New("task has no rusult")
	ErrNoSecondResult = errors.New("task has no second rusult")
)

type Task struct {
	Func reflect.Value
	Args []reflect.Value
	out  []reflect.Value

	done     bool
	result   []any
	callback func(*Task)
}

func (t *Task) AddDoneCallback(c func(*Task)) *Task {
	t.callback = c
	if t.done {
		c(t)
	}
	return t
}

func (t *Task) Done() bool {
	return t.done
}

func (t *Task) Result() []any {
	return t.result
}

func (t *Task) Default(args ...any) *Task {
	l := len(args)
	t.Args = make([]reflect.Value, l)
	for i := 0; i < l; i++ {
		t.Args[i] = reflect.ValueOf(args[i])
	}
	return t
}

func (t *Task) Reflect(args reflect.Value) *Task {
	l := args.Len()
	t.Args = make([]reflect.Value, l)
	for i := 0; i < l; i++ {
		t.Args[i] = args.Index(i)
	}
	return t
}

func (t *Task) Copy(args ...any) *Task {
	return (&Task{Func: t.Func}).Default(args...)
}

func (t *Task) Run() *Task {
	t.out = t.Func.Call(t.Args)
	t.result = make([]any, len(t.out))
	for i, v := range t.out {
		t.result[i] = v.Interface()
	}
	t.done = true
	if t.callback != nil {
		t.callback(t)
	}
	return t
}

func (t *Task) Delay(seconds float64) *Task {
	time.Sleep(time.Duration(1000*seconds) * time.Millisecond)
	return t.Run()
}

func (t *Task) First() any {
	if len(t.result) < 1 {
		panic(ErrNoResult)
	}
	return t.result[0]
}

func (t *Task) Second() any {
	if len(t.result) < 2 {
		panic(ErrNoSecondResult)
	}
	return t.result[1]
}

func (t *Task) Bool() bool {
	return t.First().(bool)
}

func (t *Task) Error() error {
	if e := t.First(); e != nil {
		return e.(error)
	}
	return nil
}

func (t *Task) Check() bool {
	return t.Error() == nil
}

func (t *Task) Bool2() bool {
	return t.Second().(bool)
}

func (t *Task) Error2() error {
	if e := t.Second(); e != nil {
		return e.(error)
	}
	return nil
}

func (t *Task) Check2() bool {
	return t.Error2() == nil
}

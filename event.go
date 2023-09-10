package asyncio

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Drelf2018/TypeGo/Pool"
	"github.com/Drelf2020/utils"
)

var eventPool = Pool.New(&Event{})

type Event struct {
	cmd     string
	data    any
	env     map[string]any
	aborted bool
	ch      chan any
}

func (e *Event) New() {
	e.env = make(map[string]any)
	e.aborted = false
}

func (e *Event) Set(x ...any) {
	e.cmd = fmt.Sprintf("%v", x[0])
	e.data = x[1]
}

func (e *Event) Reset() {
	clear(e.env)
	e.aborted = false
}

func (e *Event) Cmd() string {
	return e.cmd
}

func (e Event) String() string {
	return fmt.Sprintf("Event(%v, %v, %v)", e.cmd, e.data, e.env)
}

func (e *Event) Data(x any) error {
	return utils.CopyAny(x, e.data)
}

func (e *Event) Store(name string, value any) {
	e.env[name] = value
}

func (e *Event) Get(name string, x any, _default any) error {
	if y, ok := e.env[name]; ok {
		return utils.CopyAny(x, y)
	}
	return utils.CopyAny(x, _default)
}

func (e *Event) Abort() {
	e.aborted = true
	if e.ch != nil {
		e.ch <- struct{}{}
	}
}

func WithData[T any](handles ...func(*Event, T)) func(*Event) {
	return func(e *Event) {
		var t T
		utils.PanicErr(e.Data(&t))
		for _, h := range handles {
			if e.aborted {
				break
			}
			h(e, t)
		}
	}
}

func OnlyData[T any](handle func(T)) func(*Event) {
	return func(e *Event) {
		var t T
		utils.PanicErr(e.Data(&t))
		handle(t)
	}
}

type chain []func(*Event)

func (c chain) start(e *Event) {
	defer eventPool.Put(e)
	for i := 0; i < len(c) && !e.aborted; i++ {
		c[i](e)
	}
}

type chains []chain

func (cs chains) call(cmd, data any) {
	for _, c := range utils.NotNilSlice(cs) {
		go c.start(eventPool.Get(cmd, data))
	}
}

type model[K comparable] struct {
	chains map[K]chains
	match  func(cmd, key K) bool
}

func (m *model[K]) run(cmd K, data any) {
	Map(m.chains, func(key K, cs chains) {
		if m.match(cmd, key) {
			cs.call(cmd, data)
		}
	})
}

type AsyncEvent[K comparable] map[string]model[K]

const (
	MIN int    = -2147483648
	ALL string = "__ALL__"
)

func (a AsyncEvent[K]) Register(name string, match func(cmd, key K) bool) {
	if match == nil {
		return
	}
	a[name] = model[K]{make(map[K]chains), match}
}

func (a AsyncEvent[K]) On(name string, cmd K, handles ...func(*Event)) func() {
	m, ok := a[name]
	if !ok {
		switch name {
		case "command":
			a.Register(name, func(cmd, key K) bool { return cmd == key })
		case "regexp":
			if v, ok := any(a).(AsyncEvent[string]); ok {
				v.Register(name, func(cmd, key string) bool {
					matched, err := regexp.MatchString(key, cmd)
					return err == nil && matched
				})
				break
			}
			fallthrough
		default:
			panic("You should register \"" + name + "\" first.")
		}
		m = a[name]
	}
	m.chains[cmd] = append(m.chains[cmd], handles)
	l := len(m.chains[cmd]) - 1
	return func() { m.chains[cmd][l] = nil }
}

func (a AsyncEvent[K]) OnCommand(cmd K, handles ...func(*Event)) func() {
	return a.On("command", cmd, handles...)
}

func (a AsyncEvent[string]) OnRegexp(pattern string, handles ...func(*Event)) func() {
	return a.On("regexp", pattern, handles...)
}

func (a AsyncEvent[K]) OnAll(handles ...func(*Event)) func() {
	switch v := any(a).(type) {
	case AsyncEvent[string]:
		return v.On("command", ALL, handles...)
	case AsyncEvent[int]:
		return v.On("command", MIN, handles...)
	default:
		return nil
	}
}

func (a AsyncEvent[K]) Dispatch(cmd K, data any) {
	Map(a, func(s string, m model[K]) { m.run(cmd, data) })
	switch v := any(a).(type) {
	case AsyncEvent[string]:
		if any(cmd).(string) != ALL {
			v.Dispatch(ALL, data)
		}
	case AsyncEvent[int]:
		if any(cmd).(int) != MIN {
			v.Dispatch(MIN, data)
		}
	}
}

func Heartbeat(initdead, keepalive float64, f func(*Event)) {
	time.Sleep(time.Duration(initdead) * time.Second)

	ticker := time.NewTicker(time.Duration(keepalive) * time.Second)
	defer ticker.Stop()

	e := eventPool.Get("Heartbeat", 0)
	e.ch = make(chan any)
	defer eventPool.Put(e)

	do := func() {
		f(e)
		e.data = e.data.(int) + 1
	}

	go do()
	for {
		select {
		case <-ticker.C:
			go do()
		case <-e.ch:
			close(e.ch)
			return
		}
	}
}

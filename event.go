package asyncio

import (
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
}

func (e *Event) New() {
	e.env = make(map[string]any)
	e.aborted = false
}

func (e *Event) Set(x ...any) {
	e.cmd = x[0].(string)
	e.data = x[1]
}

func (e *Event) Reset() {
	clear(e.env)
	e.aborted = false
}

func (e *Event) Cmd() string {
	return e.cmd
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

func (cs chains) call(cmd string, data any) {
	for _, c := range utils.NotNilSlice(cs) {
		go c.start(eventPool.Get(cmd, data))
	}
}

type model struct {
	chains map[string]chains
	match  func(cmd, key string) bool
}

func (m *model) run(cmd string, data any) {
	Map(m.chains, func(key string, cs chains) {
		if m.match(cmd, key) {
			cs.call(cmd, data)
		}
	})
}

type AsyncEvent map[string]model

const ALL = "__ALL__"

func (a AsyncEvent) Register(name string, match func(cmd, key string) bool) {
	if match == nil {
		return
	}
	a[name] = model{make(map[string]chains), match}
}

func (a AsyncEvent) On(name, cmd string, handles ...func(*Event)) func() {
	m, ok := a[name]
	if !ok {
		switch name {
		case "command":
			a.Register(name, func(cmd, key string) bool { return cmd == key })
		case "regexp":
			a.Register(name, func(cmd, key string) bool {
				matched, err := regexp.MatchString(key, cmd)
				return err == nil && matched
			})
		default:
			panic("You should register \"" + name + "\" first.")
		}
		m = a[name]
	}
	m.chains[cmd] = append(m.chains[cmd], handles)
	l := len(m.chains[cmd]) - 1
	return func() { m.chains[cmd][l] = nil }
}

func (a AsyncEvent) OnCommand(cmd string, handles ...func(*Event)) func() {
	return a.On("command", cmd, handles...)
}

func (a AsyncEvent) OnRegexp(pattern string, handles ...func(*Event)) func() {
	return a.On("regexp", pattern, handles...)
}

func (a AsyncEvent) OnAll(handles ...func(*Event)) func() {
	return a.On("command", ALL, handles...)
}

func (a AsyncEvent) Dispatch(cmd string, data any) {
	Map(a, func(s string, m model) { m.run(cmd, data) })
	if cmd != ALL {
		a.Dispatch(ALL, data)
	}
}

func Heartbeat(initdead, keepalive float64, f func(*Event)) {
	time.Sleep(time.Duration(initdead) * time.Second)

	ticker := time.NewTicker(time.Duration(keepalive) * time.Second)
	defer ticker.Stop()

	e := eventPool.Get("Heartbeat", 0)
	defer eventPool.Put(e)

	do := func() {
		f(e)
		e.data = e.data.(int) + 1
	}

	go do()
	for range ticker.C {
		if e.aborted {
			break
		}
		go do()
	}
}

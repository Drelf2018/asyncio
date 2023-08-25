package asyncio

import (
	"regexp"
	"time"
)

type Event struct {
	cmd     string
	data    any
	m       map[string]any
	aborted bool
}

func (e *Event) Cmd() string {
	return e.cmd
}

func (e *Event) Data() any {
	return e.data
}

func (e *Event) Set(name string, value any) {
	e.m[name] = value
}

func (e *Event) Get(name string) any {
	return e.m[name]
}

func (e *Event) Abort() {
	e.aborted = true
}

type chain []func(*Event)

func (c chain) run(e *Event) {
	for i := 0; i < len(c) && !e.aborted; i++ {
		c[i](e)
	}
}

type chains []chain

func (cs chains) run(cmd string, data any) {
	for _, chain := range cs {
		if chain == nil {
			continue
		}
		go chain.run(&Event{cmd, data, make(map[string]any), false})
	}
}

type model struct {
	cm     map[string]chains
	choose func(cmd, key string) bool
}

func (m *model) run(cmd string, data any) {
	Map(m.cm, func(key string, cs chains) {
		if m.choose(cmd, key) {
			cs.run(cmd, data)
		}
	})
}

type AsyncEvent map[string]model

func (a AsyncEvent) Register(name string, choose func(cmd, key string) bool) {
	a[name] = model{make(map[string]chains), choose}
}

func (a AsyncEvent) Registered(name string) bool {
	_, ok := a[name]
	return ok
}

func (a AsyncEvent) On(name, cmd string, handles ...func(*Event)) func() {
	if !a.Registered(name) {
		switch name {
		case "command":
			a.Register(name, nil)
		case "regexp":
			a.Register(name, func(cmd, key string) bool {
				matched, err := regexp.MatchString(key, cmd)
				return err == nil && matched
			})
		default:
			panic("Should register \"" + name + "\" first")
		}
	}
	a[name].cm[cmd] = append(a[name].cm[cmd], handles)
	l := len(a[name].cm[cmd]) - 1
	return func() { a[name].cm[cmd][l] = nil }
}

func (a AsyncEvent) OnCommand(cmd string, handles ...func(*Event)) (delete func()) {
	return a.On("command", cmd, handles...)
}

func (a AsyncEvent) OnRegexp(pattern string, handles ...func(*Event)) (delete func()) {
	return a.On("regexp", pattern, handles...)
}

func (a AsyncEvent) All(handles ...func(*Event)) (delete func()) {
	return a.On("command", "__ALL__", handles...)
}

func (a AsyncEvent) Dispatch(cmd string, data any) {
	Map(a, func(s string, m model) {
		if s != "command" {
			m.run(cmd, data)
			return
		}
		if cs, ok := m.cm[cmd]; ok {
			cs.run(cmd, data)
		}
	})

	if cmd != "__ALL__" {
		a.Dispatch("__ALL__", data)
	}
}

func (a AsyncEvent) Heartbeat(initdead, keepalive int, f func(stop func())) {
	time.Sleep(time.Duration(initdead) * time.Second)
	ticker := time.NewTicker(time.Duration(keepalive) * time.Second)
	stopChan := make(chan any)
	for {
		select {
		case <-ticker.C:
			go f(func() { stopChan <- struct{}{} })
		case <-stopChan:
			ticker.Stop()
			return
		}
	}
}

package asyncio

import (
	"fmt"
)

type promise struct {
	EventLoop

	main    *Task
	then    []*Task
	catch   func(error)
	finally func()
}

func (p *promise) Call() {
	if p.finally != nil {
		defer p.finally()
	}
	if p.catch != nil {
		defer func() {
			if r := recover(); r != nil {
				p.catch(fmt.Errorf("%v", r))
			}
		}()
	}
	out := p.main.Run().out
	for _, then := range p.then {
		if len(then.Args) == 0 {
			then.Args = out
		}
		out = then.Run().out
	}
}

func (p *promise) Go() {
	go p.Call()
}

func (p *promise) ThenTask(task *Task) *promise {
	p.then = append(p.then, task)
	return p
}

func (p *promise) Then(fn any, args ...any) *promise {
	return p.ThenTask(p.CreateTask(fn, args...))
}

func (p *promise) Catch(fn func(error)) *promise {
	p.catch = fn
	return p
}

func (p *promise) Finally(fn func()) *promise {
	p.finally = fn
	return p
}

func PromiseTask(task *Task) *promise {
	return &promise{
		EventLoop: NewEventLoop(),
		main:      task,
	}
}

func Promise(fn any, args ...any) (p *promise) {
	p = PromiseTask(nil)
	p.main = p.CreateTask(fn, args...)
	return
}

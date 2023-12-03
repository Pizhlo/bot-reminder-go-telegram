package controller

import (
	"fmt"
	"regexp"
	"sync"

	"gopkg.in/telebot.v3"
)

type _rxhandle struct {
	origin string
	rx     *regexp.Regexp
	handle telebot.HandlerFunc
}

type Regex struct {
	handlers map[string]_rxhandle
	def      telebot.HandlerFunc
	lock     *sync.RWMutex
}

func NewRegex(def telebot.HandlerFunc) *Regex {
	p := &Regex{
		handlers: map[string]_rxhandle{},
		lock:     &sync.RWMutex{},
	}

	if def != nil {
		p.def = def
	} else {
		p.def = func(ctx telebot.Context) error {
			return nil
		}
	}

	return p
}

func (p *Regex) Handle(endpoint string, fn telebot.HandlerFunc) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.handlers[endpoint]; ok {
		panic(fmt.Errorf("handler already exists: %s", endpoint))
	}

	rx := regexp.MustCompile(endpoint)

	p.handlers[endpoint] = _rxhandle{
		origin: endpoint,
		rx:     rx,
		handle: fn,
	}
}

func (p *Regex) Handler() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		p.lock.RLock()
		defer p.lock.RUnlock()

		for _, h := range p.handlers {
			result := h.rx.FindAllStringSubmatch(ctx.Text(), -1)
			if len(result) > 0 {
				return h.handle(ctx)
			}
		}

		return p.def(ctx)
	}
}

package pool

import (
	"spexs"
	"container/list"
)

type Fifo struct {
	token chan int
	list  *list.List
}

func NewFifo() *Fifo {
	p := &Fifo{}
	p.token = make(chan int, 1)
	p.list = list.New()
	p.token <- 1
	return p
}

func (p *Fifo) Take() (*spexs.Pattern, bool) {
	<-p.token
	if p.list.Len() == 0 {
		p.token <- 1
		return nil, false
	}
	tmp := p.list.Front()
	p.list.Remove(tmp)
	p.token <- 1
	return tmp.Value.(*spexs.Pattern), true
}

func (p *Fifo) Put(pat *spexs.Pattern) {
	<-p.token
	p.list.PushBack(pat)
	p.token <- 1
}

func (p *Fifo) Len() int {
	return p.list.Len()
}
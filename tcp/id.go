package tcp

import (
	"container/list"
	"sync"
)

type idgen struct {
	used   map[int]*int
	unused *list.List
	cur    int
	mu     *sync.Mutex
}

func newIDGen() *idgen {
	return &idgen{
		mu: &sync.Mutex{},
	}
}

func (i *idgen) Get() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	if i.unused.Len() <= 0 {
		i.cur++
		elem := new(int)
		*elem = i.cur
		i.used[*elem] = elem
		return *elem
	}
	e := i.unused.Front()
	elem := e.Value.(*int)
	i.used[*elem] = elem
	i.unused.Remove(e)
	return *elem
}

func (i *idgen) Put(id int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v, ok := i.used[id]
	if ok {
		i.unused.PushBack(v)
		delete(i.used, id)
	}
}

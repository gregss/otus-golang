package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	if lc.queue.Len() == lc.capacity {
		ci := lc.queue.Back().Value.(cacheItem)
		delete(lc.items, ci.key)
		lc.queue.Remove(lc.queue.Back())
	}

	ci := cacheItem{key: key, value: value}
	if li := lc.items[key]; li != nil {
		lc.queue.MoveToFront(li)
		li.Value = ci
		return true
	}

	li := lc.queue.PushFront(ci)
	lc.items[ci.key] = li
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	li := lc.items[key]
	if li == nil {
		return nil, false
	}
	lc.queue.MoveToFront(li)

	return li.Value.(cacheItem).value, li != nil
}

func (lc *lruCache) Clear() {
	for {
		front := lc.queue.Back()
		if front == nil {
			break
		}
		ci := lc.queue.Back().Value.(cacheItem)
		delete(lc.items, ci.key)
		lc.queue.Remove(lc.queue.Back())
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

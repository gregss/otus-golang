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
	li, ok := lc.items[key]
	if !ok || li == nil {
		return nil, false
	}
	lc.queue.MoveToFront(li)

	if li.Value != nil {
		return li.Value.(cacheItem).value, true
	}

	return nil, false
}

func (lc *lruCache) Clear() {
	for {
		front := lc.queue.Back()
		if front == nil {
			break
		}
		ci, ok := lc.queue.Back().Value.(cacheItem)
		if !ok {
			return // выходим, но никому не сообщаем что работа не выполнена
		}
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

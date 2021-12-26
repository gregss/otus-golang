package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{Value: v, Next: l.front, Prev: nil}
	if l.front != nil {
		l.front.Prev = li
	}
	if l.back == nil {
		l.back = li
	}
	l.front = li
	l.len++

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.len == 0 {
		return l.PushFront(v)
	}

	l.back = &ListItem{Value: v, Next: nil, Prev: l.back}
	l.back.Prev.Next = l.back
	l.len++

	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil {
		if i.Prev != nil {
			i.Prev.Next = nil
		}
		l.back = i.Prev
	}
	if i.Prev == nil {
		if i.Next != nil {
			i.Next.Prev = nil
		}
		l.front = i.Next
	}

	if i.Next != nil && i.Prev != nil {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}

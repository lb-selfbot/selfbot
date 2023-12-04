package utils

import "sync"

type List[T comparable] struct {
	Items []T
	Mutex *sync.Mutex
}

func NewList[T comparable]() *List[T] {
	return &List[T]{Mutex: &sync.Mutex{}}
}

func (l *List[T]) Add(item T) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	l.Items = append(l.Items, item)
}

func (l *List[T]) Remove(item T) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for i, v := range l.Items {
		if v == item {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
			return
		}
	}
}

func (l *List[T]) Contains(item T) bool {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for _, v := range l.Items {
		if v == item {
			return true
		}
	}

	return false
}

func (l *List[T]) Get(index int) T {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	return l.Items[index]
}

func (l *List[T]) Len() int {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	return len(l.Items)
}

func (l *List[T]) Clear() {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	l.Items = []T{}
}

func (l *List[T]) ForEach(f func(item T)) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for _, v := range l.Items {
		f(v)
	}
}

func (l *List[T]) Filter(f func(item T) bool) *List[T] {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	filtered := NewList[T]()

	for _, v := range l.Items {
		if f(v) {
			filtered.Add(v)
		}
	}

	return filtered
}

func (l *List[T]) Map(f func(item T) T) *List[T] {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	mapped := NewList[T]()

	for _, v := range l.Items {
		mapped.Add(f(v))
	}

	return mapped
}

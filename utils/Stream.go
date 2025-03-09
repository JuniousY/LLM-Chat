package utils

import (
	"LLM-Chat/pkg"
	"sync"
	"sync/atomic"
)

type StreamResponse[T any] struct {
	l         *sync.Mutex
	q         pkg.Deque[T]
	listening bool
	listenCh  chan bool

	closed int32 // 0:未关闭 1:已关闭
}

func NewStreamResponse[T any]() *StreamResponse[T] {
	return &StreamResponse[T]{
		l:         &sync.Mutex{},
		listening: false,
		listenCh:  make(chan bool),
		closed:    0,
	}
}

func (r *StreamResponse[T]) HasNext() bool {
	r.l.Lock()

	if r.closed == 1 && r.q.Len() == 0 {
		r.l.Unlock()
		return false
	}

	if r.q.Len() > 0 {
		r.l.Unlock()
		return true
	}

	r.listening = true
	defer func() {
		r.listening = false
	}()

	r.l.Unlock()
	return <-r.listenCh
}

func (r *StreamResponse[T]) Write(v T) {
	if atomic.LoadInt32(&r.closed) == 1 {
		return
	}

	r.l.Lock()
	r.q.PushBack(v)
	if r.q.Len() == 1 && r.listening {
		r.listenCh <- true
	}
	r.l.Unlock()

	return
}

func (r *StreamResponse[T]) Read() T {
	r.l.Lock()
	defer r.l.Unlock()
	if r.q.Len() == 0 {
		var v T
		return v
	} else {
		return r.q.PopFront()
	}
}

func (r *StreamResponse[T]) Close() {
	if !atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		return
	}

	select {
	case r.listenCh <- false:
	default:
	}
	close(r.listenCh)
}

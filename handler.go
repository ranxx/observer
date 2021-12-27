package observer

import (
	"reflect"
	"sync"
)

type handler struct {
	t        reflect.Type
	observer interface{}
	callback reflect.Value
}

func newHandler(t reflect.Type, observer interface{}, function string) *handler {
	return &handler{
		t:        t,
		observer: observer,
		callback: reflect.ValueOf(observer).MethodByName(function),
	}
}

func (h *handler) Same(b *handler) bool {
	if h.t.String() == b.t.String() {
		return true
	}
	return false
}

func (h *handler) Call(params ...reflect.Value) {
	h.callback.Call(params)
}

type handlers []*handler

type syncHandler struct {
	rwlock *sync.RWMutex
	m      map[string]handlers
}

func (s *syncHandler) Append(topic string, h *handler) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	for _, v := range s.m[topic][:] {
		if v.Same(h) {
			return
		}
	}
	s.m[topic] = append(s.m[topic], h)
}

func (s *syncHandler) Range(fc func(string, handlers)) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	for topic, h := range s.m {
		fc(topic, h)
	}
}

func (s *syncHandler) Del(topic string, h *handler) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	handlers := s.m[topic]
	for i, v := range handlers[:] {
		if !v.Same(h) {
			continue
		}
		// 0，1，2，3
		left := handlers[:i]
		right := handlers[i+1:]
		left = append(left, right...)
		if len(left) <= 0 {
			delete(s.m, topic)
			return
		}
		s.m[topic] = left
		return
	}
}

func (s *syncHandler) Get(topic string) handlers {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	return s.m[topic][:]
}

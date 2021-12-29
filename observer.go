package observer

import (
	"reflect"
	"sync"
)

const (
	event = "event"
)

// Observer 观察者 模式
type Observer interface {
	Subscriber
	Publisher
	Close()
}

// Subscriber 订阅者
//
// param observer must be a reflect.Struct
//
// observer can provide Topic or Function interface
//
// or probide event field.
//
// for example:
/*
type person struct {
	event string `topic:"pain" notice:"Say"`
}

func(p*person) Say() {
	// do sth
}
*/
type Subscriber interface {
	// 订阅, 去重
	Subscribe(observer interface{})

	// 取消订阅
	Unsubscribe(observer interface{})

	// 订阅 topic 函数, 不去重, 无法取消订阅
	SubscribeByTopicFunc(topic string, fc interface{})
}

// Publisher 发布者
type Publisher interface {
	Publish(topic string, args ...interface{})
}

// Topic 主题名
type Topic interface {
	Topic() string
}

// Function 通知函数
type Function interface {
	Function() string
}

type observer struct {
	handler *syncHandler
	wait    sync.WaitGroup
}

// NewObserver Observer
func NewObserver() Observer {
	return &observer{
		handler: &syncHandler{
			rwlock: new(sync.RWMutex),
			m:      map[string]handlers{},
		},
	}
}

func (o *observer) Subscribe(observer interface{}) {
	t, topic, fc := checkObserver(observer, event)

	o.handler.Append(topic, true, newHandler(t, observer, fc))
}

func (o *observer) Unsubscribe(observer interface{}) {
	t, topic, fc := checkObserver(observer, event)

	o.handler.Del(topic, newHandler(t, observer, fc))
}

func (o *observer) SubscribeByTopicFunc(topic string, fc interface{}) {
	t, v := checkFunc(fc)

	o.handler.Append(topic, false, newHandler(t, fc, v))
}

func (o *observer) Publish(topic string, args ...interface{}) {
	handlers := o.handler.Get(topic)
	params := []reflect.Value{}

	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}

	for _, handler := range handlers {
		o.onNotice(handler, params)
	}
}

func (o *observer) onNotice(handler *handler, params []reflect.Value) {
	o.wait.Add(1)
	go func() {
		defer o.wait.Done()
		handler.Call(params...)
	}()
}

func (o *observer) Close() {
	o.wait.Wait()
}

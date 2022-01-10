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
	Wait()
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
	//
	// 返回的func是为 取消订阅
	Subscribe(observer interface{}) func()

	// // 取消订阅
	// Unsubscribe(observer interface{})

	// 订阅 topic 函数, 不去重
	//
	// 返回的func是为 取消订阅
	SubscribeByTopicFunc(topic string, fc interface{}) func()
}

// Publisher 发布者
type Publisher interface {
	// 异步推送
	Publish(topic string, args ...interface{})

	// 异步推送,监听返回值
	//
	// @params: retfc 其参数类型和顺序与消费者保持一致, 如没有返回值可为空
	PublishWithRet(topic string, retfc interface{}, args ...interface{})

	// 同步推送
	SyncPublish(topic string, args ...interface{})

	// 同步推送,监听返回值
	//
	// @params: retfc 其参数类型和顺序与消费者保持一致, 如没有返回值可为空
	SyncPublishWithRet(topic string, retfc interface{}, args ...interface{})
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

func (o *observer) Subscribe(observer interface{}) func() {
	t, topic, fc := checkObserver(observer, event)

	h := newHandler(t, observer, fc)

	o.handler.Append(topic, true, h)
	return func() {
		o.handler.Del(topic, h)
	}
}

func (o *observer) Unsubscribe(observer interface{}) {
	t, topic, fc := checkObserver(observer, event)

	o.handler.Del(topic, newHandler(t, observer, fc))
}

func (o *observer) SubscribeByTopicFunc(topic string, fc interface{}) func() {
	t, v := checkFunc(fc)

	h := newHandler(t, fc, v)

	o.handler.Append(topic, false, h)
	return func() {
		o.handler.Del(topic, h)
	}
}

func (o *observer) Publish(topic string, args ...interface{}) {
	o.publish(false, topic, nil, args...)
}

func (o *observer) PublishWithRet(topic string, retfc interface{}, args ...interface{}) {
	o.publish(false, topic, retfc, args...)
}

func (o *observer) SyncPublish(topic string, args ...interface{}) {
	o.publish(false, topic, nil, args...)
}

func (o *observer) SyncPublishWithRet(topic string, retfc interface{}, args ...interface{}) {
	o.publish(true, topic, retfc, args...)
}

func (o *observer) publish(sync bool, topic string, retfc interface{}, args ...interface{}) {
	fc := reflect.Zero(reflect.TypeOf(0))
	if retfc != nil {
		_, fc = checkFunc(retfc)
	}

	handlers := o.handler.Get(topic)
	params := []reflect.Value{}

	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}

	for _, handler := range handlers {
		o.onNotice(sync, handler, fc, params)
	}
}

func (o *observer) onNotice(sync bool, handler *handler, fc reflect.Value, params []reflect.Value) {
	h := func() {
		defer o.wait.Done()
		ret := handler.Call(params...)
		if !fc.IsZero() {
			fc.Call(ret)
		}
	}
	o.wait.Add(1)
	if sync {
		h()
		return
	}
	go h()
}

func (o *observer) Wait() {
	o.wait.Wait()
}
